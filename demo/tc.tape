# tc v0.2.0 demo video.
# Render: make demo  ->  demo/tc-v0.2.0.mp4
#
# Per knowledge-base/marketing/terminal-demos.md:
#   - Electric Sky over Obsidian; stock themes are an anti-pattern.
#   - Full-bleed, zero chrome. No window frame, card or shadow.
#   - Fit the viewport to the content: three lines of 24 columns cannot fill
#     a 16:9 canvas, so the frame hugs them instead of framing a void.
#     1200x520 is 2.3:1, inside X's 3:1 video limit.
#   - The payoff must be in the poster frame, so the tape opens on the
#     finished state and only then replays the beat.

Output demo/.scratch.mp4

Set Shell "bash"
Set FontSize 84
Set FontFamily "JetBrains Mono"
Set Width 1200
Set Height 520
Set Padding 40
Set Theme {"name":"Electric Sky","black":"#171717","red":"#F43F5E","green":"#34D399","yellow":"#FBBF24","blue":"#38BDF8","magenta":"#818CF8","cyan":"#38BDF8","white":"#E2E8F0","brightBlack":"#475569","brightRed":"#F43F5E","brightGreen":"#34D399","brightYellow":"#FBBF24","brightBlue":"#38BDF8","brightMagenta":"#818CF8","brightCyan":"#7DD3FC","brightWhite":"#F8FAFC","background":"#171717","foreground":"#C4C4C4","selection":"#1E293B","cursor":"#38BDF8"}
Set TypingSpeed 0.04s

# Open already on the done state: this is the poster frame.
Hide
Type "clear" Enter
Type "printf '💯🔥' | tc" Enter
Sleep 400ms
Show
Sleep 2500ms

# Replay the beat so the viewer sees it happen.
Hide
Type "clear" Enter
Show
Sleep 600ms
Type "printf '💯🔥' | tc"
Sleep 600ms
Enter
Sleep 3500ms
