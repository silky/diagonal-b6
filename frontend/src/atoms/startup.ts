import { b6 } from '@/lib/b6';
import { atomWithQuery } from 'jotai-tanstack-query';
import { collectionAtom } from './app';
import { viewAtom } from './location';

export const startupQueryAtom = atomWithQuery((get) => {
    const collection = get(collectionAtom);
    const viewState = get(viewAtom);
    return {
        queryKey: ['startup', collection],
        queryFn: () => {
            const { latitude, longitude } = viewState;
            const ll =
                latitude && longitude ? `${latitude},${longitude}` : undefined;
            return b6.startup({
                z: viewState.zoom.toString(),
                ...(ll ? { ll } : {}),
                r: collection,
            });
        },
    };
});
