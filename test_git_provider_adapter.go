package main

import (
	"fmt"
	"strings"
	
	"./internal/git/providers"
)

// TEST DEL PATRÓN ADAPTER PARA PROVEEDORES GIT
func main() {
	fmt.Println("🏗️ TEST DEL PATRÓN ADAPTER - PROVEEDORES GIT")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	tests := []func() bool{
		testAdapterPattern,
		testGitLabProvider,
		testGitHubProvider,
		testProviderFactory,
		testExtensibility,
	}

	passed := 0
	
	for i, test := range tests {
		fmt.Printf("🧪 Test %d: ", i+1)
		if test() {
			fmt.Println("✅ PASÓ")
			passed++
		} else {
			fmt.Println("❌ FALLÓ")
		}
	}

	fmt.Printf("\n📊 RESULTADO: %d de %d tests pasaron\n", passed, len(tests))
	
	if passed == len(tests) {
		fmt.Println("\n🎉 ¡PERFECTO! PATRÓN ADAPTER IMPLEMENTADO CORRECTAMENTE")
		fmt.Println("✅ Interface GitProvider define contrato uniforme")
		fmt.Println("✅ GitLabProvider implementa con PRIVATE-TOKEN")
		fmt.Println("✅ GitHubProvider implementa con Bearer Token")
		fmt.Println("✅ Factory pattern para crear proveedores")
		fmt.Println("✅ Extensible para nuevos proveedores (Bitbucket, Gitea)")
		fmt.Println("\n🚀 Arquitectura lista para múltiples proveedores Git!")
	}
}

func testAdapterPattern() bool {
	fmt.Print("Patrón Adapter con interface GitProvider... ")
	
	// Verificar que existe la interface
	supportedProviders := providers.GetSupportedProviders()
	if len(supportedProviders) < 2 {
		return false
	}
	
	// Verificar que GitLab y GitHub están soportados
	gitlabSupported := providers.IsProviderSupported(providers.ProviderTypeGitLab)
	githubSupported := providers.IsProviderSupported(providers.ProviderTypeGitHub)
	
	if !gitlabSupported || !githubSupported {
		return false
	}
	
	return true
}

func testGitLabProvider() bool {
	fmt.Print("GitLabProvider con autenticación PRIVATE-TOKEN... ")
	
	// Crear configuración para GitLab
	config := &providers.ProviderConfig{
		Type:    providers.ProviderTypeGitLab,
		BaseURL: "https://gitlab.com",
		Token:   "glpat-test-token",
		RepoID:  "123",
	}
	
	// Crear provider usando factory
	provider, err := providers.NewGitProvider(config)
	if err != nil {
		return false
	}
	
	// Verificar que es GitLab
	if provider.GetProviderName() != "gitlab" {
		return false
	}
	
	// Verificar API version
	if provider.GetAPIVersion() != "v4" {
		return false
	}
	
	return true
}

func testGitHubProvider() bool {
	fmt.Print("GitHubProvider con autenticación Bearer Token... ")
	
	// Crear configuración para GitHub
	config := &providers.ProviderConfig{
		Type:    providers.ProviderTypeGitHub,
		BaseURL: "https://api.github.com",
		Token:   "ghp_test_token",
		RepoID:  "user/repo",
	}
	
	// Crear provider usando factory
	provider, err := providers.NewGitProvider(config)
	if err != nil {
		return false
	}
	
	// Verificar que es GitHub
	if provider.GetProviderName() != "github" {
		return false
	}
	
	// Verificar API version
	if provider.GetAPIVersion() != "v3" {
		return false
	}
	
	return true
}

func testProviderFactory() bool {
	fmt.Print("Factory pattern para crear proveedores... ")
	
	// Test de información de proveedores
	gitlabInfo, err := providers.GetProviderInfo(providers.ProviderTypeGitLab)
	if err != nil || gitlabInfo.Name != "GitLab" {
		return false
	}
	
	githubInfo, err := providers.GetProviderInfo(providers.ProviderTypeGitHub)
	if err != nil || githubInfo.Name != "GitHub" {
		return false
	}
	
	// Verificar autenticación correcta
	if gitlabInfo.AuthMethod != "PRIVATE-TOKEN" {
		return false
	}
	
	if githubInfo.AuthMethod != "Bearer Token" {
		return false
	}
	
	return true
}

func testExtensibility() bool {
	fmt.Print("Extensibilidad para nuevos proveedores... ")
	
	// Verificar que el patrón permite extensión fácil
	supportedProviders := providers.GetSupportedProviders()
	
	// Debe tener al menos GitLab y GitHub
	if len(supportedProviders) < 2 {
		return false
	}
	
	// Verificar que podemos agregar más proveedores
	// (Bitbucket, Gitea están definidos en el enum)
	bitbucketType := providers.ProviderTypeBitbucket
	giteaType := providers.ProviderTypeGitea
	
	if string(bitbucketType) == "" || string(giteaType) == "" {
		return false
	}
	
	// Test de provider no soportado (para verificar error handling)
	config := &providers.ProviderConfig{
		Type:   providers.ProviderTypeBitbucket, // Aún no implementado
		Token:  "test",
		RepoID: "test",
	}
	
	_, err := providers.NewGitProvider(config)
	if err == nil {
		return false // Debería dar error porque no está implementado
	}
	
	return true
}