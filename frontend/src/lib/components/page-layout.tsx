import { Divider, Flex, Box } from "@chakra-ui/react";
import { ReactNode } from "react";

interface PageLayoutProps {
  children: ReactNode;
  modal: ReactNode;
}

export const PageLayout = ({ children, modal }: PageLayoutProps) => {
  return (
      <Box>
        {children}
      </Box>
  );
};
