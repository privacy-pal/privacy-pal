export default class DirectMessage {
    id: string;
    user1: string;
    user2: string;

    constructor(user1: string, user2: string) {
        this.user1 = user1;
        this.user2 = user2;
    }
}