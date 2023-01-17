package telegram

//TODO: add emoji?

const msgHelp = `I can keep your pages.
Also I'll offer you to read them.

❗After this page will be removed from my collection, so if you want to keep it, send me this link back.

💬Commands:
/help - to show help
/rnd - to get a random page
`

const msgHello = "Hi!💜  \n\n" + msgHelp

const (
	msgUnknownCommand = "I don't know this command 🤨"
	msgNoSavedPages   = "You have no saved pages 🤷‍"
	msgSaved          = "Saved 💕"
	msgAlreadyExists  = "You already sent me this page! 🙃"
)
