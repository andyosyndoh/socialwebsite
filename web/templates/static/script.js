document.getElementById('registrationForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('regUsername').value;
    const email = document.getElementById('regEmail').value;
    const password = document.getElementById('regPassword').value;
    const messageEl = document.getElementById('registrationMessage');

    try {
        const response = await fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password })
        });

        const result = await response.json();
        
        if (response.ok) {
            messageEl.innerHTML = `
                <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative" role="alert">
                    Registration Successful!
                </div>
            `;
        } else {
            messageEl.innerHTML = `
                <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                    ${result.message || 'Registration Failed'}
                </div>
            `;
        }
    } catch (error) {
        messageEl.innerHTML = `
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                Network Error: ${error.message}
            </div>
        `;
    }
});

// Login Handler
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;
    const messageEl = document.getElementById('loginMessage');

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password })
        });

        const result = await response.json();
        
        if (response.ok) {
            messageEl.innerHTML = `
                <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative" role="alert">
                    Login Successful! Welcome, ${result.username}
                </div>
            `;
        } else {
            messageEl.innerHTML = `
                <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                    ${result.message || 'Login Failed'}
                </div>
            `;
        }
    } catch (error) {
        messageEl.innerHTML = `
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                Network Error: ${error.message}
            </div>
        `;
    }
});

// Create Post Handler
document.getElementById('createPostForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const title = document.getElementById('postTitle').value;
    const content = document.getElementById('postContent').value;
    const categoriesInput = document.getElementById('postCategories').value;
    const messageEl = document.getElementById('createPostMessage');

    // Parse categories
    const categories = categoriesInput.split(',').map(cat => ({
        id: null,  // This will be set by the backend
        name: cat.trim()
    }));

    try {
        const response = await fetch('/posts/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                title, 
                content, 
                categories 
            })
        });

        const result = await response.json();
        
        if (response.ok) {
            messageEl.innerHTML = `
                <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative" role="alert">
                    Post Created Successfully! Post ID: ${result.post_id}
                </div>
            `;
        } else {
            messageEl.innerHTML = `
                <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                    ${result.message || 'Post Creation Failed'}
                </div>
            `;
        }
    } catch (error) {
        messageEl.innerHTML = `
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                Network Error: ${error.message}
            </div>
        `;
    }
});

// Retrieve Posts Handler
document.getElementById('retrievePostsForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const categoryId = document.getElementById('categoryId').value;
    const containerEl = document.getElementById('postsContainer');

    try {
        const response = await fetch(`/posts/category?category_id=${categoryId}`);
        const posts = await response.json();
        
        if (response.ok) {
            if (posts.length === 0) {
                containerEl.innerHTML = `
                    <div class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative" role="alert">
                        No posts found for this category.
                    </div>
                `;
            } else {
                const postsHTML = posts.map(post => `
                    <div class="bg-white shadow rounded-lg p-4 mb-4">
                        <h3 class="text-xl font-bold mb-2">${post.title}</h3>
                        <p class="text-gray-700 mb-2">${post.content}</p>
                        <div class="flex justify-between text-sm text-gray-500">
                            <span>Likes: ${post.likes}</span>
                            <span>Dislikes: ${post.dislikes}</span>
                        </div>
                    </div>
                `).join('');

                containerEl.innerHTML = postsHTML;
            }
        } else {
            containerEl.innerHTML = `
                <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                    Failed to retrieve posts
                </div>
            `;
        }
    } catch (error) {
        containerEl.innerHTML = `
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
                Network Error: ${error.message}
            </div>
        `;
    }
});
// Check if the user is logged in
function isUserLoggedIn() {
    // Replace this with your actual login status check
    return localStorage.getItem('isLoggedIn') === 'true';
  }
  
  // Update the UI based on login status
  function updateUI() {
    const authSection = document.getElementById('authSection');
    const authStatus = document.getElementById('authStatus');
    const mainContent = document.getElementById('mainContent');
    const loggedInUser = document.getElementById('loggedInUser');
  
    if (isUserLoggedIn()) {
      // User is logged in
      authSection.classList.add('hidden');
      authStatus.classList.remove('hidden');
      mainContent.classList.remove('hidden');
  
      // Display the logged-in user's name
      loggedInUser.textContent = localStorage.getItem('userName');
    } else {
      // User is not logged in
      authSection.classList.remove('hidden');
      authStatus.classList.add('hidden');
      mainContent.classList.add('hidden');
    }
  }
  
  // Call the updateUI function when the page loads
  window.addEventListener('DOMContentLoaded', updateUI);
  
  // Example login and registration event handlers
  document.getElementById('loginForm').addEventListener('submit', (e) => {
    e.preventDefault();
    // Perform login logic here
    localStorage.setItem('isLoggedIn', 'true');
    localStorage.setItem('userName', 'John Doe');
    updateUI();
  });
  
  document.getElementById('registrationForm').addEventListener('submit', (e) => {
    e.preventDefault();
    // Perform registration logic here
    localStorage.setItem('isLoggedIn', 'true');
    localStorage.setItem('userName', 'Jane Doe');
    updateUI();
  });