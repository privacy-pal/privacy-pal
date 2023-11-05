
export default class Message {
    id: string;
    userID: string;
    content: string;
    timestamp: Date;
    chatID: string;

    constructor(userID: string, content: string, timestamp: Date, chatID: string) {
        this.userID = userID;
        this.content = content;
        this.timestamp = timestamp;
        this.chatID = chatID;
    }
}