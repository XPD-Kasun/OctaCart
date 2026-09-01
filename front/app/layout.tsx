import { AdminShell } from "@/components/layout/Shell";
import { cn } from "@/lib/utils";
import "@/styles/style.css";
import { IBM_Plex_Sans } from "next/font/google";

const ibmPlexSans = IBM_Plex_Sans({ subsets: ['latin'], variable: '--font-sans' });


export default function Layout({ children }) {

    return (
        <html lang="en" className={cn("font-sans", ibmPlexSans.variable)}>
            <head>
                <meta charSet="UTF-8" />
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <title>OctaCard Merchant Panel</title>
            </head>
            <body>
                <AdminShell defaultRightPanelOpen={true} rightPanelOpen={true}>{children}</AdminShell>
            </body>
        </html>
    )

}