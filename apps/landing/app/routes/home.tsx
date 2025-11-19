import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "o11y - Home" },
    { name: "description", content: "Welcome to o11y!" },
  ];
}

export default function Home() {
  return <Welcome />;
}
