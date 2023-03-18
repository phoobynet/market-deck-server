const {
  protocol,
  host,
  port,
} = window.location

export const baseUrl = `${protocol}://${host}:${port}`
