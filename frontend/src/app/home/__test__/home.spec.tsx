import { render, waitFor } from "@testing-library/react";
import { VStack } from "@chakra-ui/react";
import { Home } from "../_components/home";
import { TimelinePostCard } from "../_components/timeline/timeline-post-card";
import { Post } from "@/lib/models/post";

const mockPosts: Post[] =
  [
    {
      id: "123",
      user_id: "789",
      text: "test text",
      created_at: "2024-01-01",
    },
    {
      id: "456",
      user_id: "789",
      text: "test text2",
      created_at: "2024-01-01",
    },
  ];

jest.mock("../_components/timeline/timeline-feed", () => ({
  TimelineFeed: () => {
    return (
      <VStack spacing={4} align="stretch">
        {mockPosts.map((post) => (
          <TimelinePostCard key={post.id} post={post} />
        ))}
      </VStack>
    );
  },
}));

jest.mock("next/navigation", () => ({
  useRouter() {
    return {
      push: jest.fn(),
    };
  },
}));

describe("Home Tests", () => {
  test("Rendering Home should success", () => {
    waitFor(() => {
      render(<Home />);
    });
  });
});
