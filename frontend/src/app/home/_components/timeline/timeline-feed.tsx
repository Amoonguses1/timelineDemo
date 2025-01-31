import { VStack, Text } from "@chakra-ui/react";
import { TimelinePostCard } from "./timeline-post-card";
import { Post } from "@/lib/models/post";

export const TimelineFeed = async () => {
  const posts: Post[] = [
    {
      id: "123",
      user_id: "789",
      text: "test text",
      created_at: "2024-01-01",
    },
  ];

  if (!posts || posts.length === 0) {
    return <Text>No posts found.</Text>;
  }

  return (
    <VStack spacing={4} align="stretch">
      {posts.map((post) => (
        <TimelinePostCard key={post.id} post={post} />
      ))}
    </VStack>
  );
};
