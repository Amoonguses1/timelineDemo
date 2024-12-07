"use client";

import { Box, Textarea, Button } from "@chakra-ui/react";
import { useState } from "react";

export const PostForm = () => {
  const [content, setContent] = useState("");

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    console.log("Submitted content:", content);
    setContent("");
  };

  return (
    <Box as="form" onSubmit={handleSubmit} maxW="lg" mx="auto" p={4}>
      <Textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="What's on your mind?"
        mb={3}
        minH="100px"
      />
      <Button type="submit" colorScheme="blue">
        Post
      </Button>
    </Box>
  );
};
