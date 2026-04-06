function generateRandomColors(count: number, minHueDiff = 30) {
    const colors: { h: number; s: number; l: number }[] = []
    for (let i = 0; i < count; i++) {
        let hue: number,
            attempts = 0
        const { h: previousHue } = colors?.[colors.length - 1] ?? {}
        do {
            hue = Math.random() * 360
            attempts++
        } while (attempts < 100 && previousHue !== undefined && Math.abs(previousHue - hue) < minHueDiff)

        colors.push({ h: hue, s: 70, l: 75 })
    }
    return colors.map((c) => `hsl(${c.h}, ${c.s}%, ${c.l}%)`)
}

export default generateRandomColors
