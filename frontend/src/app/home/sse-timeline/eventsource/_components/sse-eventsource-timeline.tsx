"use client";

import { VStack, Box } from "@chakra-ui/react";
import { TimelinePostCard } from "../../../_components/timeline/timeline-post-card";
import { useSSEEventSourceTimelineFeed } from "./use-sse-eventsource-timeline";

export const SSEEventSourceTimelineFeed = () => {
  const { posts, errorMessage } = useSSEEventSourceTimelineFeed();

  if (errorMessage) {
    // Handling errors that cannot be caught by error.tsx from asynchronous processing.
    return <Box>{errorMessage}</Box>;
  } else if (posts.length === 0) {
    return <Box>Post not found.</Box>;
  } else {
    return (
      <VStack spacing={4} align="stretch">
        {posts.map((post) => (
          <TimelinePostCard key={`sse-eventsource-${post.id}`} post={post} />
        ))}
      </VStack>
    );
  }
};
