{ pkgs
, b6-go-packages
, pyproject-nix
}:
let
  # We have to write the version into here from the output of the go
  # binary.
  pyproject-file = (pkgs.runCommand "make-pyproject" { } ''
    substitute ${./../python/pyproject.toml.template} $out \
      --subst-var-by VERSION ''$(${b6-go-packages.go-executables.b6-api}/bin/b6-api --pip-version)
  '');

  pythonProject = pyproject-nix.lib.project.loadPyproject {
    projectRoot = ./../python;
    pyproject = pkgs.lib.importTOML pyproject-file;
  };

  renderedPyProject = python: pythonProject.renderers.buildPythonPackage {
    inherit python;
  };

  # A couple of hacks necessary to build the proto files and the API.
  pythonHacks = ''
    # Set the pyproject to be the one we computed via our b6 binary.
    cat ${pyproject-file} > pyproject.toml

    # Bring in the necessary proto files and the Makefile
    cp -r ${./../proto} ./proto
    cat ${./../Makefile} > some-Makefile

    # Hack: Run the b6-api command outside of the Makefile, using the
    # Nix version of the binary.
    ${b6-go-packages.go-executables.b6-api}/bin/b6-api --functions \
      | python diagonal_b6/generate_api.py > diagonal_b6/api_generated.py

    # Hack: Make the directory structure that the Makefile expects,
    # then move things to where we want them
    mkdir python
    mkdir python/diagonal_b6

    make proto-python -f some-Makefile

    # Cleanup
    mv python/diagonal_b6/* ./diagonal_b6
    rm -rf python/diagonal_b6
  '';

  wheel = pkgs.stdenv.mkDerivation {
    name = "b6-wheel";
    src = ../.;

    nativeBuildInputs = [
      pkgs.python312.pkgs.grpcio-tools
    ];

    buildInputs = with pkgs; [
      (python312.withPackages (ps: [ ps.build ps.setuptools ]))
    ];

    patchPhase = ''
      cd python
      ${pythonHacks}
    '';

    buildPhase = ''
      python -m build -n
    '';

    installPhase = ''
      mv dist $out
    '';
  };

  # Note: This library can only be included on the `python` that is provided
  # here.
  b6-py = python: python.pkgs.buildPythonPackage ((renderedPyProject python) // {
    nativeBuildInputs = [
      python.pkgs.grpcio-tools
    ];

    preBuild = ''
      ${pythonHacks}
    '';

    pythonImportsCheck = [ "diagonal_b6" ];
  });
in
{
  inherit b6-py wheel;
}
