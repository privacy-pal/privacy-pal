import { Locator } from './locator';
import { UpdateData } from 'firebase-admin/firestore';

export type HandleAccessFunc = (dataSubjectId: string, locator: Locator, obj: any) => Record<string, any>;
export type HandleDeletionFunc = (dataSubjectId: string) => {
    nodesToTraverse: Locator[],
    deleteNode: boolean,
    updateData?: UpdateData<any>
};