export default function Layout({ children }: { children: React.ReactNode }) {
    return (
        <div>
            <p>Header</p>
            {children}
            <p>Footer</p>
        </div>
    )
}
