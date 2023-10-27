export default class Message {
    id: string;
    userID: string;
    content: string;
    timestamp: Date;

    constructor(userID: string, content: string, timestamp: Date) {
        this.userID = userID;
        this.content = content;
        this.timestamp = timestamp;
    }
}