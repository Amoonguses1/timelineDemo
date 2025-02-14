import { useRef, useState, useEffect } from "react";
import { ERROR_MESSAGES } from "@/lib/constants/error-messages";
import { Post } from "@/lib/models/post";

export interface UseSSETimelineFeedReturn {
  errorMessage: string | null;
  posts: Post[];
}

export type SSETimelineEventResponse =
  | TimelineAccessedResponse
  | PostCreatedResponse
  | PostDeletedResponse;

export enum SSETimelineEventType {
  TimelineAccessed = "TimelineAccessed",
  PostCreated = "PostCreated",
  PostDeleted = "PostDeleted",
}

interface TimelineAccessedResponse {
  EventType: SSETimelineEventType.TimelineAccessed;
  Posts: Post[];
}

interface PostCreatedResponse {
  EventType: SSETimelineEventType.PostCreated;
  Posts: Post[];
}

interface PostDeletedResponse {
  EventType: SSETimelineEventType.PostDeleted;
  Posts: Post[];
}

export const useSSEEventSourceTimelineFeed = (): UseSSETimelineFeedReturn => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);
  const url: string = `${process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL}/api/${process.env.NEXT_PUBLIC_USER_ID}/sse`;

  useEffect(() => {
    if (eventSourceRef.current) return;

    const eventSource = new EventSource(url);

    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      try {
        const newPosts: SSETimelineEventResponse = JSON.parse(event.data);
        switch (newPosts.EventType) {
          case SSETimelineEventType.TimelineAccessed:
          case SSETimelineEventType.PostCreated:
            if (!newPosts.Posts || newPosts.Posts.length === 0) {
              return;
            }
            setPosts((prevPosts) => [...newPosts.Posts, ...prevPosts]);
            break;
          case SSETimelineEventType.PostDeleted:
            // TODO: https://github.com/okuda-seminar/Twitter-Clone/issues/540
            // - Implement timeline post deletion with SSE event handling.
            return;
        }
      } catch (err) {
        setErrorMessage(ERROR_MESSAGES.INVALID_DATA);
      }
    };

    eventSource.onerror = () => {
      setErrorMessage(ERROR_MESSAGES.SERVER_ERROR);
      eventSource.close();
      eventSourceRef.current = null;
    };

    return () => {
      eventSource.close();
      eventSourceRef.current = null;
    };
  }, [url]);

  return {
    errorMessage,
    posts,
  };
};
