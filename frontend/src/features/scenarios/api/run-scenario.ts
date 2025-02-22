import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";

import { getEvaluate } from "@/api/evaluate";
import { World, useWorldStore } from "@/stores/worlds";
import { getWorldFeatureId } from "@/utils/world";

import { useChangesStore } from "../stores/changes";
import { useComparisonsStore } from "../stores/comparisons";
import { useTabsStore } from "../stores/tabs";
import { ChangeSpec } from "../types/change";

export const useScenario = (
	origin: World,
	target: World,
	change: ChangeSpec,
) => {
	const actions = useChangesStore((state) => state.actions);
	const tabActions = useTabsStore((state) => state.actions);
	const { setFeatureId, setTiles } = useWorldStore((state) => state.actions);
	const { add: addComparison } = useComparisonsStore((state) => state.actions);
	const query = useQuery({
		enabled: false,
		queryKey: ["scenario", origin.id, target.id, JSON.stringify(change)],
		queryFn: () => {
			if (!change.changeFunction?.id) {
				return Promise.reject("no change function defined");
			}
			if (!change.features || change.features.length === 0) {
				return Promise.reject("no features defined");
			}

			return getEvaluate({
				root: origin.featureId,
				request: {
					call: {
						function: {
							symbol: "add-world-with-change",
						},
						args: [
							{
								literal: {
									featureIDValue: target.featureId,
								},
							},
							{
								call: {
									function: {
										call: {
											function: {
												symbol: "evaluate-feature",
											},
											args: [
												{
													literal: {
														featureIDValue: change.changeFunction.id,
													},
												},
											],
										},
									},
									args: [
										{
											literal: {
												collectionValue: {
													keys: change.features.map((_, i) => {
														return {
															intValue: i,
														};
													}),
													values: change.features.map((f) => {
														return {
															featureIDValue: f.id,
														};
													}),
												},
											},
										},
									],
								},
							},
						],
					},
				},
			});
		},
	});

	useEffect(() => {
		if (query.isSuccess) {
			actions.setCreate(target.id, true);
			tabActions.setPersist(target.id, true);
			const featureId = getWorldFeatureId({
				...origin.featureId,
				namespace: `${origin.featureId.namespace}/scenario`,
				value: target.featureId.value,
			});

			setFeatureId(target.id, featureId);
			setTiles(target.id, target.id);

			if (change.analysis) {
				addComparison({
					id: `${origin.id}-${target.id}`,
					baseline: origin,
					scenarios: [
						{
							id: target.id,
							featureId: featureId,
							tiles: target.id,
						},
					],
					analysis: change.analysis,
				});
			}
		}
	}, [query.isSuccess, actions, target.id, setFeatureId]);

	return query;
};
