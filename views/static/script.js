document.addEventListener("DOMContentLoaded", function () {
  const modal = document.getElementById("auth-modal");
  const loginForm = document.getElementById("login-form");
  const signupForm = document.getElementById("signup-form");
  const showSignup = document.getElementById("show-signup");
  const showLogin = document.getElementById("show-login");
  const loginButton = document.getElementById("login");
  const closeModal = document.getElementById("close-modal");
  const modalContent = document.querySelector("#auth-modal .form-container");
  const headerActions = document.querySelector(".header-actions");
  const createbutton = document.getElementById("create-post-btn");
  const logoutButton = document.getElementById("logout");
  const postsContainer = document.getElementById("posts-container");

  // Fetch posts from the backend
  function fetchAndRenderPosts() {
    fetch("/posts")
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error("Failed to fetch posts");
      })
      .then((posts) => {
        // Clear the posts container
        postsContainer.innerHTML = "";

        // Loop through the posts and render each one
        posts.forEach((post) => {
          const postElement = createPostElement(post);
          postsContainer.appendChild(postElement);
        });
      })
      .catch((error) => {
        console.error("Error fetching posts:", error);
      });
  }

  // Helper function to create a post element
  function createPostElement(post) {
    const postCard = document.createElement("div");
    postCard.classList.add("post-card", "card");

    postCard.innerHTML = `
          <div class="post-header">
              <p class="post-author">${post.username || "Anonymous"}</p>
              <p class="post-time">${new Date(
                post.timestamp
              ).toLocaleString()}</p>
          </div>
          <div class="post-content">
              <p>${post.content}</p>
          </div>
          <div class="post-actions">
              <button class="like-btn" data-post-id="${post.id}">Like (${
      post.likes
    })</button>
              <button class="dislike-btn" data-post-id="${post.id}">Dislike (${
      post.dislikes
    })</button>
              <button class="comment-btn" data-post-id="${
                post.id
              }">Comment</button>
          </div>
      `;

    // Add event listeners for the like, dislike, and comment buttons
    const likeButton = postCard.querySelector(".like-btn");
    const dislikeButton = postCard.querySelector(".dislike-btn");
    const commentButton = postCard.querySelector(".comment-btn");

    likeButton.addEventListener("click", () => handleLike(post.id));
    dislikeButton.addEventListener("click", () => handleDislike(post.id));
    commentButton.addEventListener("click", () => handleComment(post.id));

    return postCard;
  }

  // Handlers for like, dislike, and comment actions
  function handleLike(postId) {
    fetch(`/posts/${postId}/like`, { method: "POST" })
      .then(fetchAndRenderPosts)
      .catch((err) => console.error("Error liking post:", err));
  }

  function handleDislike(postId) {
    fetch(`/posts/${postId}/dislike`, { method: "POST" })
      .then(fetchAndRenderPosts)
      .catch((err) => console.error("Error disliking post:", err));
  }

  function handleComment(postId) {
    alert(`Redirect to comment functionality for post ${postId}`);
  }

  // Fetch and render posts on page load
  fetchAndRenderPosts();

  cookie();

  logoutButton.addEventListener("click", function () {
    fetch("/logout", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
    })
      .then((response) => {
        if (response.ok) {
          window.location.reload();
        } else {
          console.error("Logout failed");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  });

  // Show modal when login/signup button is clicked
  loginButton.addEventListener("click", function () {
    modal.classList.remove("hidden");
  });

  // Close modal
  closeModal.addEventListener("click", function () {
    modal.classList.add("hidden");
  });

  // Switch to signup form
  showSignup.addEventListener("click", function (e) {
    e.preventDefault();
    loginForm.classList.remove("active");
    loginForm.classList.add("hidden");
    signupForm.classList.remove("hidden");
    signupForm.classList.add("active");
  });

  // Switch to login form
  showLogin.addEventListener("click", function (e) {
    e.preventDefault();
    signupForm.classList.remove("active");
    signupForm.classList.add("hidden");
    loginForm.classList.remove("hidden");
    loginForm.classList.add("active");
  });

  // Prevent closing modal when clicking inside it
  modal.addEventListener("click", function (e) {
    if (e.target === modal) {
      modal.classList.add("hidden");
    }
  });

  // Handle signup form submission
  signupForm.addEventListener("submit", function (e) {
    e.preventDefault(); // Prevent default form submission behavior

    const username = signupForm.querySelector("#signup-username").value.trim();
    const email = signupForm.querySelector("#signup-email").value.trim();
    const password = signupForm.querySelector("#signup-password").value.trim();

    const data = { username, email, password };

    fetch("/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
      .then((response) => {
        if (response.ok) {
          signupForm.classList.remove("active");
          signupForm.classList.add("hidden");
          loginForm.classList.remove("hidden");
          loginForm.classList.add("active");
          displayMessage(modalContent, "");
        } else {
          throw new Error("Registration failed");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        displayMessage(modalContent, "Unsuccessful. Try again.");
      });
  });

  // Handle login form submission
  loginForm.addEventListener("submit", function (e) {
    e.preventDefault(); // Prevent default form submission behavior

    const username = loginForm.querySelector("#login-username").value.trim();
    const password = loginForm.querySelector("#login-password").value.trim();

    const data = { username, password };

    fetch("/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
      .then((response) => {
        if (response.ok) {
          return response.json(); // Assuming the server sends some JSON data on success
        } else {
          throw new Error("Login failed");
        }
      })
      .then((data) => {
        console.log("Logged in:", data);

        // Hide login/signup button
        loginButton.style.display = "none";

        // Disable login/signup modals
        modal.classList.add("hidden");

        // Optionally, you can add user info to the page (e.g., "Welcome, [username]")
        cookie();
      })
      .catch((error) => {
        console.error("Error:", error);
        displayMessage(modalContent, "Not successful. Try again.");
      });
  });

  /**
   * Utility function to display a message within the modal
   * @param {Element} container - The container to display the message in
   * @param {string} message - The message to display
   */
  function displayMessage(container, message) {
    let messageElement = container.querySelector(".modal-message");

    if (!messageElement) {
      messageElement = document.createElement("p");
      messageElement.className = "modal-message";
      messageElement.style.color = "var(--accent-color)";
      messageElement.style.marginTop = "1rem";
      container.appendChild(messageElement);
    }

    messageElement.textContent = message;
  }

  function cookie() {
    fetch("/profile")
      .then((response) => {
        if (response.ok) {
          return response.json();
        }
        throw new Error("Not logged in");
      })
      .then((user) => {
        console.log("User details:", user);

        // Show 'Create Post' and 'Logout' buttons, hide 'Login' button
        document.getElementById("login").style.display = "none";
        document.getElementById("create-post-btn").style.display =
          "inline-block";
        document.getElementById("logout").style.display = "inline-block";

        // Create a user greeting and display user's name
        const userGreeting = document.createElement("p");
        userGreeting.textContent = `Welcome, ${user.username || "User"}!`;
        userGreeting.style.color = "var(--accent-color)";

        // Assuming you have a headerActions element where the greeting should be appended
        const headerActions = document.querySelector(".header-actions");
        headerActions.appendChild(userGreeting);

        // Optionally, display the username in the 'login' button as well
        document.getElementById(
          "login"
        ).textContent = `Welcome, ${user.username}`;
      })
      .catch((error) => {
        console.error("Error:", error);

        // Show 'Login' button when user is not logged in, hide 'Create Post' and 'Logout' buttons
        document.getElementById("login").style.display = "inline-block";
        createbutton.style.display = "none";
        document.getElementById("logout").style.display = "none";
      });
  }
});
