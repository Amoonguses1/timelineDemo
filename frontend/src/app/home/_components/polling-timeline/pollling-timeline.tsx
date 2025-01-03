"use client";
import React, { useState, useEffect } from "react";
import { VStack } from "@chakra-ui/react";
import { TimelinePostCard } from "../timeline/timeline-post-card";
import { getInitialTimeline } from "@/lib/actions/get_initial_timeline";
import { pollFollowingPosts } from "@/lib/actions/poll_following_post";
import { Post } from "@/lib/models/post";

export const PollingTimelineFeed = () => {
  const [posts, setPosts] = useState<Post[]>([]);
  const { data: initialData, error: initialError } = getInitialTimeline();
  const { data: pollingData, error: pollingError } = pollFollowingPosts();
  const initialPosts = initialData;
  const newPosts = pollingData;

  useEffect(() => {
    if (initialPosts?.length) {
      setPosts(initialPosts);
    }
  }, [initialData]);

  useEffect(() => {
    if (newPosts?.length) {
      setPosts((prevPosts) => [...newPosts, ...prevPosts]);
    }
  }, [pollingData]);

  if (initialError) {
    return (
      <div>
        <p>failed to initial fetch</p>
        <p>error message: {initialError.message}</p>
      </div>
    );
  } else if (pollingError) {
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
          <TimelinePostCard key={post.id} post={post} />
        ))}
      </VStack>
    );
  }
};
