import Box from '@mui/material/Box';
import { WallMessage } from './types';
import MessageBox from './MessageBox';

export default function MessageStream(props: {messages: WallMessage[]}) {
    const { messages } = props;

    const orderedMessages = messages.reverse()
    return (
        <Box>
        {
            orderedMessages.map((message) => {
                return (
                    <MessageBox key={message.id} message={message}></MessageBox>
                )
            })
        }
        </Box>
    )
}