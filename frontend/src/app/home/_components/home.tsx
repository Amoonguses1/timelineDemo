import React from "react";
import { TimelineFeed } from "./timeline/timeline-feed";
import { Divider, Box, Flex, VStack } from "@chakra-ui/react";
import { PollingTimelineFeed } from "./polling-timeline/pollling-timeline";
import { getInitialTimeline } from "@/lib/actions/get_initial_timeline";

export const Home = async () => {
  const initialPosts = await getInitialTimeline();
  return (
    <Flex width="100%" height="100vh">
      <Box flex="1 1 33%">
        <VStack>
          <Box fontSize="lg">Normal timeline feed</Box>
          <Divider borderColor="white" />
          <Box width="100%">
            <TimelineFeed />
          </Box>
        </VStack>
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 33%">
        <VStack>
          <Box fontSize="lg">Polling timeline feed</Box>
          <Divider borderColor="white" />
          <Box width="100%">
            <PollingTimelineFeed initialPosts={initialPosts} />
          </Box>
        </VStack>
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 33%">
        <VStack>
          <Box fontSize="lg">SSE timeline feed</Box>
          <Divider borderColor="white" />
          <Box width="100%">
            <PollingTimelineFeed initialPosts={initialPosts} />
          </Box>
        </VStack>
      </Box>
    </Flex>
  );
};
