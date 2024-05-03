import { useOutlinerContext } from '@/lib/context/outliner';
import colors from '@/tokens/colors.json';
import { HistogramBarLineProto, SwatchLineProto } from '@/types/generated/ui';
import { scaleOrdinal } from '@visx/scale';
import { interpolateRgbBasis } from 'd3-interpolate';
import { useCallback, useEffect, useMemo } from 'react';
import { match } from 'ts-pattern';
import { Histogram } from '../system/Histogram';

const colorInterpolator = interpolateRgbBasis([
    '#fff',
    colors.amber[20],
    colors.violet[80],
]);

type HistogramData = {
    index: number;
    label: string;
    count: number;
};

export const HistogramAdaptor = ({
    type,
    bars,
    swatches,
}: {
    type: 'swatch' | 'histogram';
    bars?: HistogramBarLineProto[];
    swatches?: SwatchLineProto[];
}) => {
    const { outliner, setHistogramColorScale, setHistogramBucket } =
        useOutlinerContext();
    const scale = outliner.histogram?.colorScale;

    const data = useMemo(() => {
        return match(type)
            .with(
                'histogram',
                () =>
                    bars?.flatMap((bar) => {
                        return {
                            index: bar.index ?? 0,
                            label: bar.range?.value ?? '',
                            count: bar.value,
                        };
                    }) ?? []
            )
            .with(
                'swatch',
                () =>
                    swatches?.flatMap((swatch) => {
                        return {
                            index: swatch.index ?? 0,
                            label: swatch.label?.value ?? '',
                            /* Swatches do not have a count. Should be null but setting it to 0 
                            for now to avoid type errors. */
                            count: 0,
                        };
                    }) ?? []
            )
            .exhaustive();
    }, [type, bars, swatches]);

    useEffect(() => {
        const scale = scaleOrdinal({
            domain: data.map((d) => `${d.index}`),
            range: data.map((_, i) => colorInterpolator(i / data.length)),
        });
        setHistogramColorScale(scale);
    }, [data]);

    const handleSelect = useCallback((d: HistogramData | null) => {
        setHistogramBucket(d?.index.toString());
    }, []);

    const selected = useMemo(() => {
        const selected = outliner?.histogram?.selected;
        if (!selected) return null;
        return data.find((d) => d.index.toString() === selected);
    }, [outliner.histogram?.selected, data]);

    return (
        <Histogram
            type={type}
            data={data}
            label={(d) => d.label}
            bucket={(d) => d.index.toString()}
            value={(d) => d.count}
            color={(d) => (scale ? scale(`${d.index}`) : '#fff')}
            onSelect={handleSelect}
            selected={selected}
            selectable
        />
    );
};