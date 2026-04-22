package handler

import (
	"fmt"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/spf13/cobra"
)

// MarketplaceHandler handles marketplace subcommands for all primitive types.
type MarketplaceHandler struct {
	marketplace *service.MarketplaceService
	install     *service.InstallService
}

// NewMarketplaceHandler creates a MarketplaceHandler.
func NewMarketplaceHandler(mp *service.MarketplaceService, install *service.InstallService) *MarketplaceHandler {
	return &MarketplaceHandler{marketplace: mp, install: install}
}

// List returns a RunE handler for "cap <type> marketplace list".
func (h *MarketplaceHandler) List(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		primitives, err := h.marketplace.List(primitiveType)
		if err != nil {
			return err
		}
		if len(primitives) == 0 {
			cmd.Println("No primitives available. Run 'cap update' to sync repositories.")
			return nil
		}

		total := len(primitives)
		start := (page - 1) * size
		if start >= total {
			return fmt.Errorf("page %d out of range", page)
		}
		end := start + size
		if end > total {
			end = total
		}

		cmd.Printf("%-35s  %s\n", "NAME", "VERSION")
		for _, p := range primitives[start:end] {
			cmd.Printf("%-35s  %s\n", p.Vendor+"/"+p.Name, p.Version)
		}
		pages := (total + size - 1) / size
		cmd.Printf("\nСтраница %d из %d\n", page, pages)
		return nil
	}
}

// Search returns a RunE handler for "cap <type> marketplace search <regex>".
func (h *MarketplaceHandler) Search(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		results, err := h.marketplace.Search(primitiveType, args[0])
		if err != nil {
			return err
		}
		if len(results) == 0 {
			cmd.Println("No matches found.")
			return nil
		}
		cmd.Printf("%-35s  %-10s  %s\n", "NAME", "VERSION", "DESCRIPTION")
		for _, p := range results {
			cmd.Printf("%-35s  %-10s  %s\n", p.Vendor+"/"+p.Name, p.Version, p.Description)
		}
		return nil
	}
}

// Info returns a RunE handler for "cap <type> marketplace info <name> [<version>]".
func (h *MarketplaceHandler) Info(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := ""
		if len(args) > 1 {
			version = args[1]
		}

		results, err := h.marketplace.Info(primitiveType, name, version)
		if err != nil {
			return err
		}

		// results are sorted latest-first by MarketplaceService.Info.
		p := results[0]
		cmd.Printf("Name:        %s\n", p.Name)
		cmd.Printf("Version:     %s\n", p.Version)
		cmd.Printf("Vendor:      %s\n", p.Vendor)
		cmd.Printf("Namespace:   cap:%s:%s\n", p.Vendor, p.Name)
		cmd.Printf("Description: %s\n", p.Description)

		if len(results) > 1 {
			cmd.Println("\nAvailable versions:")
			seen := map[string]bool{}
			for _, r := range results {
				if !seen[r.Version] {
					cmd.Printf("  - %s\n", r.Version)
					seen[r.Version] = true
				}
			}
		}
		return nil
	}
}

// Add returns a RunE handler for "cap <type> marketplace add <name> [<version>]".
func (h *MarketplaceHandler) Add(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := ""
		if len(args) > 1 {
			version = args[1]
		}

		p, err := h.install.Add(primitiveType, name, version)
		if err != nil {
			return err
		}
		cmd.Printf("Installed: cap:%s:%s  version %s\n", p.Vendor, p.Name, p.Version)
		return nil
	}
}
