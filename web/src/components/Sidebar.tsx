import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Button,
  Divider,
  Flex,
  Heading,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  useStyleConfig,
  VStack,
} from "@chakra-ui/react";
import { routes } from "@/constants/routes";
import { useBudget } from "@/data/queries/budget";
import { useUser } from "@/context/AuthContext";
import { generatePath, Link } from "react-router-dom";
import { useLogout } from "@/data/queries/user";

const links = [
  { name: "home", pathname: routes.budget.index },
  { name: "budget", pathname: routes.budget.budget },
  { name: "accounts", pathname: routes.budget.accounts },
  { name: "transactions", pathname: routes.budget.transactions },
  { name: "settings", pathname: routes.budget.settings },
];

const Sidebar = () => {
  const user = useUser();
  const styles = useStyleConfig("Sidebar");
  const { budget } = useBudget();

  const { mutate: logout, isLoading: isLoggingOut } = useLogout();

  return (
    <Flex
      __css={styles}
      p={4}
      h="100vh"
      position="sticky"
      top="0"
      w={56}
      boxShadow="md"
      flexShrink="0"
    >
      <Flex direction="column">
        <Heading size="md">beans</Heading>
        <Divider my={3} />
        <Button
          to={routes.budget.noneSelected}
          as={Link}
          leftIcon={<ArrowBackIcon />}
          size="xs"
          w="full"
        >
          {budget.name}
        </Button>
        <VStack align="flex-start" mt={6}>
          {links.map(({ name, pathname }) => (
            <Button
              key={pathname}
              to={generatePath(pathname, { budget: budget.id })}
              as={Link}
              size="sm"
              w="full"
              justifyContent="flex-start"
              variant="ghost"
            >
              {name}
            </Button>
          ))}
        </VStack>
      </Flex>
      <Flex direction="column">
        <Divider my={3} />
        <Menu>
          <MenuButton as={Button} variant="ghost" size="sm" textAlign="left">
            {user?.username}
          </MenuButton>
          <MenuList>
            <MenuItem onClick={() => logout()} disabled={isLoggingOut}>
              Log out
            </MenuItem>
          </MenuList>
        </Menu>
      </Flex>
    </Flex>
  );
};

export default Sidebar;
