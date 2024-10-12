import {
    Tag,
    Users,
    Settings,
    Bookmark,
    SquarePen,
    LayoutGrid,
    LucideIcon
  } from "lucide-react";
  
  type Submenu = {
    href: string;
    label: string;
    active?: boolean;
  };
  
  type Menu = {
    href: string;
    label: string;
    active?: boolean;
    icon: LucideIcon;
    submenus?: Submenu[];
  };
  
  type Group = {
    groupLabel: string;
    menus: Menu[];
  };
  
  export function getMenuList(pathname: string): Group[] {
    return [
      {
        groupLabel: "",
        menus: [
          {
            href: "/dashboard",
            label: "Dashboard",
            icon: LayoutGrid,
            submenus: []
          }
        ]
      },
      {
        groupLabel: "Contents",
        menus: [
          {
            href: "",
            label: "Agent",
            icon: SquarePen,
            submenus: [
              {
                href: "/dashboard/agents",
                label: "All Agents"
              },
              {
                href: "/dashboard/agents/new",
                label: "New Agent"
              }
            ]
          },
          // {
          //   href: "/dashboard/categories",
          //   label: "Categories",
          //   icon: Bookmark
          // },
          // {
          //   href: "/dashboard/tags",
          //   label: "Tags",
          //   icon: Tag
          // }
        ]
      },
      {
        groupLabel: "Settings",
        menus: [
          {
            href: "/dashboard/users",
            label: "Users",
            icon: Users
          },
          // {
          //   href: "/dashboard/account",
          //   label: "Account",
          //   icon: Settings
          // }
        ]
      }
    ];
  }