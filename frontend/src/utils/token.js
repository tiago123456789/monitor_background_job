const getValueInPayloud = (key, token) => {
    const tokenSplited = token.split(".")
    let payload = atob(tokenSplited[1])
    payload = JSON.parse(payload)
    return payload[key] || null
}

export default {
    getValueInPayloud
}