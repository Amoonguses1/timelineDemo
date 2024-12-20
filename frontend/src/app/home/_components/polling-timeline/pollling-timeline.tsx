"use client";
import { VStack } from "@chakra-ui/react";
import { TimelinePostCard } from "../timeline/timeline-post-card";
import { getInitialTimeline } from "@/lib/actions/get_initial_timeline.tsx";

export const PollingTimelineFeed = () => {
  const { data: initialData, error: initialError } = getInitialTimeline();

  const posts = initialData;

  if (initialError) {
    return (
      <div>
        <p>failed to initial fetch</p>
        <p>error message: {initialError.message}</p>
      </div>
    );
  } else if (!posts) {
    return <div>post not found</div>;
  } else {
    return (
      <VStack spacing={4} align="stretch">
        {posts.map((post) => (
          <TimelinePostCard key={post.id} post={post} />
        ))}
      </VStack>
    );
  }
};
