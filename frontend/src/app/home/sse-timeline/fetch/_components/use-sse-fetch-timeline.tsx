"use client";
import { useState, useEffect } from "react";
import { ERROR_MESSAGES } from "@/lib/constants/error-messages";
import { Post } from "@/lib/models/post";
import {
  SSETimelineEventResponse,
  SSETimelineEventType,
  UseSSETimelineFeedReturn,
} from "../../eventsource/_components/use-sse-eventsource-timeline";

export const useSSEFetchTimelineFeed = (): UseSSETimelineFeedReturn => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    const url = `${process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL}/api/${process.env.NEXT_PUBLIC_USER_ID}/sse`;
    let abortController = new AbortController();

    const fetchData = async () => {
      const res = await fetch(url);
      const reader = res.body?.getReader()!;
      const decoder = new TextDecoder();

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        if (!value) continue;

        const lines = decoder.decode(value);

        try {
          const firstLine = lines
            .split("\n")
            .find((line) => line.startsWith("data: "));

          if (!firstLine) continue;

          const jsonData = firstLine.replace(/^data: /, "");
          const newPosts = JSON.parse(jsonData) as SSETimelineEventResponse;

          switch (newPosts.EventType) {
            case SSETimelineEventType.TimelineAccessed:
            case SSETimelineEventType.PostCreated:
              setPosts((prevPosts) => [...newPosts.Posts, ...prevPosts]);
              break;
            case SSETimelineEventType.PostDeleted:
              // TODO: https://github.com/okuda-seminar/Twitter-Clone/issues/540
              // - Implement timeline post deletion with SSE event handling.
              break;
          }
        } catch (parseError) {
          console.error("Error parsing SSE data:", parseError);
          setErrorMessage(ERROR_MESSAGES.INVALID_DATA);
        }
      }
    };
    fetchData();
    return () => {
      abortController.abort();
      setPosts([]);
      setErrorMessage(null);
    };
  }, []);

  return {
    errorMessage,
    posts,
  };
};
