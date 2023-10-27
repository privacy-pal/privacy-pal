import Message from "./message"

export default class GroupChat {
    id: string;
    users: string[];
    messages: Message[];
    owner: string;

    constructor(owner: string, users: string[]) {
        this.owner = owner;
        this.users = users;
        this.messages = [];
    }
}