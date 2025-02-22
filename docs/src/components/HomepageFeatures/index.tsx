import clsx from "clsx";
import Link from "@docusaurus/Link";
import Heading from "@theme/Heading";
import styles from "./styles.module.css";

type FeatureItem = {
	title: string;
	Svg: React.ComponentType<React.ComponentProps<"svg">>;
	description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
	{
		title: "API",
		description: (
			<>The b6 gRPC API; most commonly consumed through the Python library.</>
		),
		link: "/docs/api",
	},
	{
		title: "Backend",
		description: (
			<>
				The b6 backend, written in Go, providing the webserver, gRPC API server,
				map tiles, custom rendering, and various ingest/post-processing tools.
			</>
		),
		link: "/docs/backend",
	},
	{
		title: "Frontend",
		description: <>The b6 frontend, written in React, providing the map UI.</>,
		link: "/docs/frontend",
	},
	{
		title: "Contributing",
		description: (
			<>How to contribute to the codebase; building from source, etc.</>
		),
		link: "/docs/contributing",
	},
	{
		title: "Quirks",
		description: <>Quirks, bugs, and known issues.</>,
		link: "/docs/quirks",
	},
];

function Feature({ title, Svg, description, link, label }: FeatureItem) {
	return (
		<div className={clsx("col col--4")}>
			<div className="text--center padding-horiz--md">
				<Heading as="h3">
					<Link className="button button--secondary button--lg" to={link}>
						{title}
					</Link>
				</Heading>
				<p>{description}</p>
			</div>
		</div>
	);
}

export default function HomepageFeatures(): JSX.Element {
	return (
		<section className={styles.features}>
			<div className="container">
				<div className="row">
					{FeatureList.map((props, idx) => (
						<Feature key={idx} {...props} />
					))}
				</div>
			</div>
		</section>
	);
}
