import { useAppContext } from '@/lib/context/app';
import { OutlinerProvider, OutlinerStore } from '@/lib/context/outliner';
import {
    DndContext,
    MouseSensor,
    PointerSensor,
    TouchSensor,
    useDraggable,
    useDroppable,
    useSensor,
    useSensors,
} from '@dnd-kit/core';
import { restrictToWindowEdges } from '@dnd-kit/modifiers';
import { AnimatePresence, motion } from 'framer-motion';
import { PropsWithChildren } from 'react';
import { twMerge } from 'tailwind-merge';
import { StackAdapter } from './adapters/StackAdapter';

export const OutlinersLayer = ({ mapId }: { mapId: string }) => {
    const {
        draggableOutliners,
        dockedOutliners,
        setActiveOutliner,
        setFixedOutliner,
        moveOutliner,
    } = useAppContext();
    const pointerSensor = useSensor(PointerSensor, {
        activationConstraint: {
            distance: 5,
        },
    });
    const mouseSensor = useSensor(MouseSensor);
    const touchSensor = useSensor(TouchSensor);

    const sensors = useSensors(pointerSensor, mouseSensor, touchSensor);

    return (
        <>
            <div className="absolute top-16 left-2 flex flex-col gap-1">
                {dockedOutliners.map((outliner) => {
                    return (
                        <OutlinerProvider key={outliner.id} outliner={outliner}>
                            <StackAdapter />
                        </OutlinerProvider>
                    );
                })}
            </div>
            <DndContext
                modifiers={[restrictToWindowEdges]}
                sensors={sensors}
                onDragStart={({ active }) => {
                    setActiveOutliner(active.id as string, true);
                    setFixedOutliner(active.id as string);
                }}
                onDragEnd={({ active, delta }) => {
                    moveOutliner(active.id as string, delta.x, delta.y);
                    setActiveOutliner(active.id as string, false);
                }}
            >
                <Droppable mapId={mapId}>
                    <AnimatePresence>
                        {draggableOutliners.map((outliner) => {
                            return (
                                <DraggableOutliner
                                    key={outliner.id}
                                    outliner={outliner}
                                />
                            );
                        })}
                    </AnimatePresence>
                </Droppable>
            </DndContext>
        </>
    );
};

const Droppable = ({
    children,
    mapId,
}: PropsWithChildren & { mapId: string }) => {
    const { isOver, setNodeRef } = useDroppable({
        id: `droppable-${mapId}`,
    });
    const style = {
        color: isOver ? 'green' : undefined,
    };
    return (
        <div ref={setNodeRef} style={style}>
            {children}
        </div>
    );
};

const DraggableOutliner = ({
    outliner,
}: PropsWithChildren & {
    outliner: OutlinerStore;
}) => {
    const active = outliner.active;
    const { attributes, transform, setNodeRef, listeners } = useDraggable({
        id: outliner.id,
    });

    const style = {
        transform: `${
            transform
                ? `translate3d(${transform.x}px, ${transform.y}px, 0)`
                : ''
        }`,
    };

    const variants = {
        hidden: {
            opacity: 0,
            scale: 0,
        },
        visible: {
            opacity: 1,
            scale: 1,
        },
    };

    return (
        <div
            id={outliner.id}
            ref={setNodeRef}
            style={{
                ...style,
                top: outliner.properties.coordinates?.y + 4,
                left: outliner.properties.coordinates?.x + 4,
                position: 'absolute',
            }}
            className={twMerge(
                active && 'ring-2 ring-ultramarine-50 ring-opacity-40'
            )}
            {...listeners}
            {...attributes}
        >
            <motion.div
                variants={variants}
                initial="hidden"
                animate="visible"
                exit="hidden"
                transition={{
                    duration: 0.1,
                }}
            >
                <div>
                    <OutlinerProvider outliner={outliner}>
                        <StackAdapter />
                    </OutlinerProvider>
                </div>
            </motion.div>
        </div>
    );
};