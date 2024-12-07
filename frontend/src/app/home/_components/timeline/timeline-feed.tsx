import { VStack, Text } from "@chakra-ui/react";
import { getCollectionOfPostsBySpecificUserAndUsersTheyFollow, GetCollectionOfPostsBySpecificUserAndUsersTheyFollowResponse } from "@/lib/actions/get-collection-of-posts-by-specific-user-and-users-they-follow";
import { TimelinePostCard } from "./timeline-post-card";

export const TimelineFeed = async () => {
  // const posts = await getCollectionOfPostsBySpecificUserAndUsersTheyFollow({
  //   user_id: `${process.env.NEXT_PUBLIC_USER_ID}`,
  // });

  const posts: GetCollectionOfPostsBySpecificUserAndUsersTheyFollowResponse =
    [
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
