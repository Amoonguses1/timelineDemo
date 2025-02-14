import React from "react";
import { TimelineFeed } from "@/app/home/_components/timeline/timeline-feed";
import { Divider, Box, Flex, VStack } from "@chakra-ui/react";
import { SSEEventSourceTimelineFeed } from "./sse-eventsource-timeline";
import Link from "next/link";

export const SSEEventSourceTimelineHome = async () => {
  return (
    <Flex width="100%" height="100vh">
      <Box flex="1 1 50%">
        <Link href={"/home"}>
          <VStack>
            <Box fontSize="lg">timeline feed</Box>
            <Divider borderColor="white" />
            <Box width="100%">
              <TimelineFeed />
            </Box>
          </VStack>
        </Link>
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 50%">
        <VStack>
          <Box fontSize="lg">SSE Eventsource timeline feed</Box>
          <Divider borderColor="white" />
          <Box width="100%">
            <SSEEventSourceTimelineFeed />
          </Box>
        </VStack>
      </Box>
    </Flex>
  );
};
