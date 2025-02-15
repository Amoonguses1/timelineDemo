"use client";
import React, { useState, useEffect } from "react";
import { VStack } from "@chakra-ui/react";
import { TimelinePostCard } from "../../_components/timeline/timeline-post-card";
import { pollFollowingPosts } from "@/lib/actions/poll_following_post";
import { Post } from "@/lib/models/post";

export type InitialPostsProps = {
  initialPosts: Post[];
};

export const PollingTimelineFeed = ({ initialPosts }: InitialPostsProps) => {
  const [posts, setPosts] = useState<Post[]>(initialPosts);
  const { data: pollingPosts, error: pollingError } = pollFollowingPosts();

  useEffect(() => {
    if (pollingPosts?.length) {
      setPosts((prevPosts) => [...pollingPosts, ...prevPosts]);
    }
  }, [pollingPosts]);

  if (pollingError) {
    return (
      <div>
        <p>failed to polling fetch</p>
        <p>error message: {pollingError.message}</p>
      </div>
    );
  } else if (!posts) {
    return <div>post not found</div>;
  } else {
    return (
      <VStack spacing={4} align="stretch">
        {posts.map((post) => (
          <TimelinePostCard key={`polling-${post.id}`} post={post} />
        ))}
      </VStack>
    );
  }
};
