import o11yLogo from "./logo.png"
import grafanaLogo from "./grafana.svg"
import prometheusLogo from "./prometheus.svg"
import { useContext, useEffect, useState } from "react";
import { appContext } from "~/context";

export function Welcome() {

  const [loggedIn, setLoggedIn ] = useState<Boolean>(false)
  const [err, setErr] = useState<string>("");
  const [user, setUser] = useState<string>("")
  const { authUrl } = useContext(appContext);


  useEffect(() => {
    async function getLoggedInUser() {
      fetch(`${authUrl}/me`, {
        credentials: 'include',
      }).then(async (resp) => {
        if (resp.ok) {
          const { email } = await resp.json();
          setUser(email);
          setLoggedIn(true);
        } else {
          setLoggedIn(false);
        }
      }).catch((err) => {
        setLoggedIn(false);
      }) 
    }
    getLoggedInUser();
  }, [loggedIn])

  
  function logIn(formData: FormData) {
    console.log("Logging in...")
    fetch(`${authUrl}/login`, {
      method: "POST",
      credentials: 'include',
      body: JSON.stringify({
        email: formData.get("email"),
        password: formData.get("password")
      })
    }).then((resp) => {
      if (!resp.ok) {
        setErr("something wrong");
      } else {
        setLoggedIn(true);
      }
    })
  }

  return (
    <main className="flex items-center justify-center pt-16 pb-4">
      <div className="flex-1 flex flex-col items-center gap-16 min-h-0">
        <header className="flex flex-col items-center gap-9">
          <div className="w-[200px] max-w-[100vw] p-4">
            <img
              src={o11yLogo}
              alt="o11y Logo"
              className="block w-full"
            />
          </div>
          <p className="text-4xl">o11y</p>
        </header>
        <div className="max-w-[300px] w-full space-y-6 px-4">
          <nav className="rounded-3xl border border-gray-200 p-6 dark:border-gray-700 space-y-4">
            {loggedIn && (
              <p className="leading-6 text-gray-700 dark:text-gray-200 text-center"> Logged in as: {user}</p>
            ) || (
              <div>
                <p className="leading-6 text-gray-700 dark:text-gray-200 text-center">Log in</p>
                <form action={logIn}>
                  <input className="bg-gray-700 text-center border-1 border-solid border-gray-600 m-1" name="email" type="email" placeholder="email"/>
                  <input className="bg-gray-700 text-center border-1 border-solid border-gray-600 m-1" name="password" type="password" placeholder="password" />
                  <button type="submit" hidden></button>
                </form>
              </div>
            )}
          </nav>
        </div>
        <div className="max-w-[300px] w-full space-y-6 px-4">
          <nav className="rounded-3xl border border-gray-200 p-6 dark:border-gray-700 space-y-4">
            <p className="leading-6 text-gray-700 dark:text-gray-200 text-center">
              Available tools
            </p>
            <ul>
              {resources.map(({ href, text, icon }) => (
                <li key={href}>
                  <a
                    className="group flex items-center gap-3 self-stretch p-3 leading-normal text-blue-700 hover:underline dark:text-blue-500"
                    href={href}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {icon}
                    {text}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </div>
    </main>
  );
}

const resources = [
  {
    href: "https://grafana.o11y.local",
    text: "Grafana",
    icon: (
      <img
        width="24"
        height="20"
        src={grafanaLogo}
      />
    ),
  },
  {
    href: "https://prometheus.o11y.local",
    text: "Prometheus",
    icon: (
      <img
        width="24"
        height="20"
        src={prometheusLogo}
      />
    ),
  }
];
