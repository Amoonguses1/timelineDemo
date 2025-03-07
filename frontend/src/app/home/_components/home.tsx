import React from "react";
import { Divider, Box, Flex, VStack } from "@chakra-ui/react";
import Link from "next/link";

export const Home = async () => {
  return (
    <Flex width="100%" height="100vh">
      <Box flex="1 1 33%">
        <Link href={"/home/polling-timeline"}>
          <VStack>
            <Box fontSize="lg">Polling timeline feed</Box>
            <Divider borderColor="white" />
          </VStack>
        </Link>
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 33%">
        <Link href={"/home/sse-timeline/eventsource"}>
          <VStack>
            <Box fontSize="lg">SSE Eventsource timeline feed</Box>
            <Divider borderColor="white" />
          </VStack>
        </Link>
      </Box>
      <Divider orientation="vertical" borderColor="white" />
      <Box flex="1 1 33%">
        <Link href={"/home/sse-timeline/fetch"}>
          <VStack>
            <Box fontSize="lg">SSE fetch timeline feed</Box>
            <Divider borderColor="white" />
          </VStack>
        </Link>
      </Box>
    </Flex>
  );
};
