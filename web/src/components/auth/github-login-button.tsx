import { GitBranch } from 'lucide-react'

interface GithubLoginButtonProps {
  label: string
}

export function GithubLoginButton({ label }: GithubLoginButtonProps) {
  return (
    <a className="auth-submit auth-submit-github" href="/api/v1/user/oauth/github">
      <GitBranch size={17} aria-hidden="true" />
      <span>{label}</span>
    </a>
  )
}
