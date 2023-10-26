import { Locator } from './locator';
import { UpdateData } from 'firebase-admin/firestore';

export interface DataNode {
    handleAccess(dataSubjectId: string, locator: Locator): Record<string, any>;
    handleDeletion(dataSubjectId: string): {
        nodesToTraverse: Locator[],
        deleteNode: boolean,
        updateData?: UpdateData<any>
    };
}
