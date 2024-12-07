import React from "react";
import { TimelineFeed } from "./timeline/timeline-feed";
import { Divider, Box,Flex } from "@chakra-ui/react";
import { PostForm } from "./post-form/post-form";

export const Home = () => {
  return (
    <Flex width="100%" height="100vh" >
      <Box flex="1 1 50%">
        <TimelineFeed />
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 50%">
        <PostForm />
      </Box>
    </Flex>
  );
};
