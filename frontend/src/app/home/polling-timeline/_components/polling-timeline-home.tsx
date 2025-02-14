import React from "react";
import { TimelineFeed } from "../../_components/timeline/timeline-feed";
import { Divider, Box, Flex, VStack } from "@chakra-ui/react";
import { PollingTimelineFeed } from "./pollling-timeline";
import { getInitialTimeline } from "@/lib/actions/get_initial_timeline";
import Link from "next/link";

export const PollingTimelineHome = async () => {
  const initialPosts = await getInitialTimeline();
  return (
    <Flex width="100%" height="100vh">
      <Box flex="1 1 33%">
        <Link href={"/home"}>
          <VStack>
            <Box fontSize="lg">Normal timeline feed</Box>
            <Divider borderColor="white" />
            <Box width="100%">
              <TimelineFeed />
            </Box>
          </VStack>
        </Link>
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
    </Flex>
  );
};
