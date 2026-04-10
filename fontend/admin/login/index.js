// Importar funciones necesarias del config.js centralizado
import { login, API_CONFIG } from '../../services/config.js';
import { Auth } from '../../services/auth.js';

document.getElementById('loginForm').addEventListener('submit', async (e) => {
  e.preventDefault();
  
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  const errorDiv = document.getElementById('loginError');
  const submitBtn = e.target.querySelector('button[type="submit"]');
  const submitText = submitBtn ? submitBtn.querySelector('span') : null;

  errorDiv.style.display = 'none';
  submitBtn.disabled = true;
  if (submitText) submitText.innerHTML = 'Iniciando sesión...';

  try {
    // Usar la función login del config.js centralizado
    const response = await login(username, password);

    // El servicio de usuarios devuelve accessToken, no token
    if (response.accessToken) {
      Auth.setToken(response.accessToken);
      Auth.setUser(response.user);
      if (response.refreshToken) Auth.setRefreshToken(response.refreshToken);
      Auth.setLastActivity();

      // Mostrar éxito
      if (submitText) submitText.innerHTML = '¡Bienvenido!';

      setTimeout(() => {
        window.location.href = '../dashboard/';
      }, 500);
    } else {
      throw new Error(response.message || 'Error al iniciar sesión');
    }
  } catch (error) {
    errorDiv.querySelector('p').textContent = error.message || '';
    errorDiv.style.display = 'block';
    submitBtn.disabled = false;
    if (submitText) submitText.innerHTML = 'Acceder';
  }
});

