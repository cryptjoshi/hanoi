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
  
  function isActive(menuPath: string, currentPath: string): boolean {
    return currentPath.startsWith(menuPath);
  }
  
  export function getMenuList(pathname: string, lng: string | undefined): Group[] {
    // Ensure lng is a string, use a default if it's undefined
    const language = typeof lng === 'string' ? lng : 'en';

    return [
      {
        groupLabel: "",
        menus: [
          {
            href: `/${language}/dashboard`,
            label: "dashboard",
            icon: LayoutGrid,
            active: isActive(`/${language}/dashboard`, pathname),
            submenus: []
          }
        ]
      },
      {
        groupLabel: "contents",
        menus: [
          {
            href: "",
            label: "agent",
            icon: SquarePen,
            active: isActive(`/${language}/dashboard/agents`, pathname),
            submenus: [
              {
                href: `/${language}/dashboard/agents`,
                label: "all_agents",
                active: isActive(`/${language}/dashboard/agents`, pathname)
              },
              {
                href: `/${language}/dashboard/agents/new`,
                label: "new_agent",
                active: pathname === `/${language}/dashboard/agents/new`
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
        groupLabel: "settings",
        menus: [
          {
            href: `/${language}/dashboard/settings`,
            label: "setting",
            icon: Settings,
            active: isActive(`/${language}/dashboard/settings`, pathname)
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
