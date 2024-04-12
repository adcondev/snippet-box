// Select all anchor tags within the navigation
const navLinks = document.querySelectorAll("nav a");

// Loop through each navigation link
for (let i = 0; i < navLinks.length; i++) {
    let link = navLinks[i]

    // If the href attribute of the link matches the current path
    if (link.getAttribute('href') == window.location.pathname) {
        // Add the "live" class to the link
        link.classList.add("live");

        // Break the loop as we've found the active link
        break;
    }
}