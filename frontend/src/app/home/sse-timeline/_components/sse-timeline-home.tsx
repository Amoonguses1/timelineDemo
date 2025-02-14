import React from "react";
import { Divider, Box, Flex, VStack } from "@chakra-ui/react";
import Link from "next/link";

export const SSETimelineHome = async () => {
  return (
    <Flex width="100%" height="100vh">
      <Box flex="1 1 33%">
        <Link href={"sse-timeline/eventsource"}>
          <VStack>
            <Box fontSize="lg">SSE EventSource timeline feed</Box>
            <Divider borderColor="white" />
          </VStack>
        </Link>
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 33%">
        <Link href={"sse-timeline/fetch"}>
          <VStack>
            <Box fontSize="lg">SSE fetch timeline feed</Box>
            <Divider borderColor="white" />
          </VStack>
        </Link>
      </Box>
    </Flex>
  );
};
