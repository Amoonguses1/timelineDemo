import { useRef, useState, useEffect } from "react";
import { ERROR_MESSAGES } from "@/lib/constants/error-messages";
import { Post } from "@/lib/models/post";

interface useSSETimelineFeedReturn {
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
  event_type: SSETimelineEventType.TimelineAccessed;
  posts: Post[];
}

interface PostCreatedResponse {
  event_type: SSETimelineEventType.PostCreated;
  posts: Post[];
}

interface PostDeletedResponse {
  event_type: SSETimelineEventType.PostDeleted;
  posts: Post[];
}

export const useSSETimelineFeed = (): useSSETimelineFeedReturn => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);
  const url: string = `${process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL}/api/users/${process.env.NEXT_PUBLIC_USER_ID}/timelines/reverse_chronological`;

  useEffect(() => {
    if (eventSourceRef.current) return;

    const eventSource = new EventSource(url);

    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      try {
        const newPosts: SSETimelineEventResponse = JSON.parse(event.data);
        switch (newPosts.event_type) {
          case SSETimelineEventType.TimelineAccessed:
          case SSETimelineEventType.PostCreated:
            if (!newPosts.posts || newPosts.posts.length === 0) {
              return;
            }
            setPosts((prevPosts) => [...newPosts.posts, ...prevPosts]);
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
