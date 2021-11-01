package bot

// RegistrationStatusInit is a constant that indicates the temporary user's registration status is in init stage
const RegistrationStatusInit = 0

// RegistrationStatusUserName is a constant that indicates the temporary user has registered username
const RegistrationStatusUserName = 1

// RegistrationStatusPhoneNumber is a constant that indicates the temporary user has registered phonenumber
const RegistrationStatusPhoneNumber = 2

// RegistrationStatusCategory is a constant that indicates the temporary user has selected user category
const RegistrationStatusCategory = 3

// MainMenuW is a constant that holds the main menu value with post job button
var MainMenuW = CreateReplyKeyboard(true, false,
	[]string{"📋 Post Job", "💼 Manage Jobs"},
	[]string{"🔔 Job Subscriptions", "⚙️ Settings"})

// MainMenuWO is a constant that holds the main menu value with out post job button
var MainMenuWO = CreateReplyKeyboard(true, false,
	[]string{"💼 Manage Jobs"},
	[]string{"🔔 Job Subscriptions", "⚙️ Settings"})

// SubscriptionModified is a constant that indicates a subscription field has been added to previously created one
const SubscriptionModified = 1

// SubscriptionNotFound is a constant that indicates the given subscription is not found
const SubscriptionNotFound = 2

// SubscriptionError is a constant that indicates an error related to modifying subscription
const SubscriptionError = 3
