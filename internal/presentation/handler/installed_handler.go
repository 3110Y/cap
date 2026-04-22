package handler

import (
	"fmt"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/spf13/cobra"
)

// InstalledHandler handles commands for primitives installed in the project.
type InstalledHandler struct {
	svc *service.InstalledService
}

// NewInstalledHandler creates an InstalledHandler.
func NewInstalledHandler(svc *service.InstalledService) *InstalledHandler {
	return &InstalledHandler{svc: svc}
}

// List returns a RunE handler for "cap <type> list".
func (h *InstalledHandler) List(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		size, _ := cmd.Flags().GetInt("size")

		primitives, err := h.svc.List(primitiveType)
		if err != nil {
			return err
		}
		if len(primitives) == 0 {
			cmd.Printf("No %s installed.\n", primitiveType)
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

		cmd.Printf("%-30s  %s\n", "NAME", "VERSION")
		for _, p := range primitives[start:end] {
			cmd.Printf("%-30s  %s\n", p.Vendor+"/"+p.Name, p.Version)
		}
		pages := (total + size - 1) / size
		cmd.Printf("\nСтраница %d из %d\n", page, pages)
		return nil
	}
}

// Info returns a RunE handler for "cap <type> info <name>".
func (h *InstalledHandler) Info(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		info, err := h.svc.Info(primitiveType, args[0])
		if err != nil {
			return err
		}

		cmd.Printf("Name:        %s\n", info.Name)
		cmd.Printf("Vendor:      %s\n", info.Vendor)
		cmd.Printf("Namespace:   cap:%s:%s\n", info.Vendor, info.Name)
		if info.Description != "" {
			cmd.Printf("Description: %s\n", info.Description)
		}

		cmd.Println("\nInstalled versions:")
		for _, v := range info.InstalledVersions {
			cmd.Printf("  - %s\n", v)
		}

		if len(info.AvailableVersions) > 0 {
			cmd.Println("\nAvailable versions (marketplace):")
			for _, v := range info.AvailableVersions {
				cmd.Printf("  - %s\n", v)
			}
		}
		return nil
	}
}

// Delete returns a RunE handler for "cap <type> del <name>".
func (h *InstalledHandler) Delete(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		if err := h.svc.Delete(primitiveType, args[0], version); err != nil {
			return err
		}
		if version != "" {
			cmd.Printf("Removed: %s  version %s\n", args[0], version)
		} else {
			cmd.Printf("Removed: %s  (all versions)\n", args[0])
		}
		return nil
	}
}

// Upgrade returns a RunE handler for "cap <type> upgrade <name>".
func (h *InstalledHandler) Upgrade(primitiveType string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		p, err := h.svc.Upgrade(primitiveType, args[0])
		if err != nil {
			return err
		}
		cmd.Printf("Upgraded: cap:%s:%s  →  %s\n", p.Vendor, p.Name, p.Version)
		return nil
	}
}

// UpgradeAll is a RunE handler for "cap upgrade".
func (h *InstalledHandler) UpgradeAll(cmd *cobra.Command, _ []string) error {
	upgraded, errs := h.svc.UpgradeAll()

	for _, p := range upgraded {
		cmd.Printf("Upgraded: cap:%s:%s  →  %s\n", p.Vendor, p.Name, p.Version)
	}
	for _, err := range errs {
		cmd.PrintErrf("Warning: %v\n", err)
	}

	if len(upgraded) == 0 && len(errs) == 0 {
		cmd.Println("All primitives are up to date.")
	} else {
		cmd.Printf("\n%d primitive(s) upgraded.\n", len(upgraded))
	}
	return nil
}
