// fix_schema.go - OpenAPI Schema Fixer
//
// This script fixes incorrect schema references in the VRChat OpenAPI specification.
// Specifically, it replaces "type: string" with "$ref: '#/components/schemas/UserID'"
// for parameters that should reference User IDs, and fixes external UserID references.
//
// The script fixes:
// 1. External schema references: "../schemas/UserID.yaml" → "#/components/schemas/UserID"
// 2. Parameters explicitly named "userId", "userIdAdmin", "confirmEmailUserId", etc.
// 3. Parameters with examples containing "usr_" (User ID pattern)
// 4. Parameters with descriptions mentioning "user ID" or "UserID"
// 5. Adds missing schema definitions to fix code generation errors
//
// The script uses pattern-based matching in either:
// - Direct parameter definitions (under parameters: sections)
// - Component parameter definitions (under components: but not in schemas:)
//
// This approach is resilient to file changes, unlike line-number-based fixes.
//
// Usage: go run fix_schema.go (from the utils directory)

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Read the openapi.yaml file from parent directory
	filename := filepath.Join("..", "openapi.yaml")
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	content := string(data)

	// Split content into lines for easier manipulation
	lines := strings.Split(content, "\n")

	// Track changes made
	changes := 0

	// First, fix external schema references to UserID.yaml
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "$ref: ../schemas/UserID.yaml" {
			indentation := strings.Repeat(" ", len(lines[i])-len(strings.TrimLeft(lines[i], " ")))
			lines[i] = indentation + "$ref: '#/components/schemas/UserID'"
			changes++
			fmt.Printf("Fixed line %d: external UserID.yaml reference\n", i+1)
		}
	}

	// Fix patterns based on parameter context instead of line numbers
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		
		// Check if this line contains a parameter definition (should end with : and not be a sub-property)
		if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") && len(line) > 1 {
			parameterName := strings.TrimSuffix(line, ":")
			
			// Skip common non-parameter names
			if parameterName == "parameters" || parameterName == "schema" || parameterName == "properties" || 
			   parameterName == "components" || parameterName == "paths" || parameterName == "deprecated" ||
			   strings.HasPrefix(parameterName, "$") || strings.HasPrefix(parameterName, "/") ||
			   parameterName == "name" || parameterName == "in" || parameterName == "required" ||
			   parameterName == "description" || parameterName == "example" || parameterName == "type" {
				continue
			}
			
			// Look for parameters that should be UserID references
			shouldBeUserID := false
			
			// Check if parameter name explicitly indicates it's a User ID
			lowerName := strings.ToLower(parameterName)
			if strings.Contains(lowerName, "userid") || 
			   strings.HasSuffix(lowerName, "userid") ||
			   lowerName == "confirmemuserid" ||
			   (strings.Contains(lowerName, "user") && strings.Contains(lowerName, "id")) {
				shouldBeUserID = true
			}
			
			// Also check for patterns in the parameter definition
			if !shouldBeUserID {
				parameterBlock := ""
				for j := i; j < i+15 && j < len(lines); j++ {
					// Stop if we hit another top-level parameter
					if j > i && strings.HasSuffix(strings.TrimSpace(lines[j]), ":") && 
					   !strings.HasPrefix(lines[j], "  ") {
						break
					}
					parameterBlock += strings.ToLower(lines[j]) + " "
				}
				
				// Check for usr_ pattern in examples AND user-related descriptions
				if strings.Contains(parameterBlock, "usr_") &&
				   (strings.Contains(parameterBlock, "user id") || 
				    strings.Contains(parameterBlock, "userid") ||
				    strings.Contains(parameterBlock, "valid user") || 
				    strings.Contains(parameterBlock, "target user")) {
					shouldBeUserID = true
				}
			}
			
			// If this should be a UserID parameter, look for its schema definition
			if shouldBeUserID {
				for j := i + 1; j < i+15 && j < len(lines); j++ {
					if strings.Contains(lines[j], "schema:") {
						// Look for type: string in the next few lines after schema
						for k := j + 1; k < j+8 && k < len(lines); k++ {
							if strings.TrimSpace(lines[k]) == "type: string" {
								// Check if this is in a valid parameter section
								isValidParameterSection := false
								inComponentsSection := false
								inSchemasSection := false
								
								for l := i - 200; l < i && l >= 0; l++ {
									if strings.Contains(lines[l], "parameters:") {
										isValidParameterSection = true
										break
									}
									if strings.Contains(lines[l], "components:") {
										inComponentsSection = true
									}
									if strings.Contains(lines[l], "  schemas:") {
										inSchemasSection = true
									}
								}
								
								// Valid if we found parameters: OR if we're in components but not in schemas
								if isValidParameterSection || (inComponentsSection && !inSchemasSection) {
									indentation := strings.Repeat(" ", len(lines[k])-len(strings.TrimLeft(lines[k], " ")))
									lines[k] = indentation + "$ref: '#/components/schemas/UserID'"
									changes++
									fmt.Printf("Fixed line %d: %s parameter schema (detected as UserID)\n", k+1, parameterName)
								}
								break
							}
						}
						break
					}
				}
			}
		}
	}

	// Fix platform_history in CurrentUser to use PlatformHistory schema reference
	for i := 0; i < len(lines)-20; i++ {
		if strings.TrimSpace(lines[i]) == "platform_history:" {
			// Check if this is in the CurrentUser schema by looking backwards
			inCurrentUserSchema := false
			for j := i - 200; j < i && j >= 0; j++ {
				if strings.Contains(lines[j], "CurrentUser:") {
					inCurrentUserSchema = true
					break
				}
			}
			
			if inCurrentUserSchema {
				// Look for the inline array definition
				if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "type: array" {
					// Find the items section and replace it
					for k := i + 1; k < i + 20 && k < len(lines); k++ {
						if strings.TrimSpace(lines[k]) == "items:" {
							// Look for the inline object definition
							if k+1 < len(lines) && strings.TrimSpace(lines[k+1]) == "type: object" {
								// Find the end of the properties block
								endOfProperties := k + 1
								for l := k + 2; l < len(lines); l++ {
									line := strings.TrimSpace(lines[l])
									// Stop when we hit the next property at the same level as platform_history
									if line != "" && !strings.HasPrefix(lines[l], "          ") && 
									   !strings.HasPrefix(lines[l], "            ") &&
									   strings.HasSuffix(line, ":") {
										endOfProperties = l - 1
										break
									}
								}
								
								// Replace the entire items block with a reference
								indentation := strings.Repeat(" ", len(lines[k])-len(strings.TrimLeft(lines[k], " ")))
								newLines := make([]string, 0, len(lines))
								newLines = append(newLines, lines[:k+1]...)
								newLines = append(newLines, indentation+"  $ref: '#/components/schemas/PlatformHistory'")
								newLines = append(newLines, lines[endOfProperties+1:]...)
								lines = newLines
								changes++
								fmt.Printf("Fixed platform_history in CurrentUser to use PlatformHistory schema reference\n")
								break
							}
						}
					}
				}
			}
		}
	}

	// Fix inline otp schemas to use Otp schema reference
	for i := 0; i < len(lines)-20; i++ {
		if strings.TrimSpace(lines[i]) == "otp:" {
			// Look for the inline array definition with items: type: object
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "type: array" {
				// Find the items section and check if it has an inline object
				for k := i + 1; k < i + 20 && k < len(lines); k++ {
					if strings.TrimSpace(lines[k]) == "items:" {
						// Look for the inline object definition
						if k+1 < len(lines) && strings.TrimSpace(lines[k+1]) == "type: object" {
							// This is an inline otp schema, find the end of the object
							endOfObject := k + 1
							for l := k + 2; l < len(lines); l++ {
								line := strings.TrimSpace(lines[l])
								// Stop when we hit the next property at the same level as otp
								if line != "" && !strings.HasPrefix(lines[l], "          ") && 
								   !strings.HasPrefix(lines[l], "            ") &&
								   strings.HasSuffix(line, ":") {
									endOfObject = l - 1
									break
								}
							}
							
							// Replace the entire items block with a reference
							indentation := strings.Repeat(" ", len(lines[k])-len(strings.TrimLeft(lines[k], " ")))
							newLines := make([]string, 0, len(lines))
							newLines = append(newLines, lines[:k+1]...)
							newLines = append(newLines, indentation+"  $ref: '#/components/schemas/Otp'")
							newLines = append(newLines, lines[endOfObject+1:]...)
							lines = newLines
							changes++
							fmt.Printf("Fixed otp inline schema to use Otp schema reference\n")
							break
						}
					}
				}
			}
		}
	}

	// Fix inline publishedListings schemas to use PublishedListing schema reference
	for i := 0; i < len(lines)-20; i++ {
		if strings.TrimSpace(lines[i]) == "publishedListings:" {
			// Look for the inline array definition with items: type: object
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "type: array" {
				// Find the items section and check if it has an inline object
				for k := i + 1; k < i + 20 && k < len(lines); k++ {
					if strings.TrimSpace(lines[k]) == "items:" {
						// Look for the inline object definition
						if k+1 < len(lines) && strings.TrimSpace(lines[k+1]) == "type: object" {
							// This is an inline publishedListings schema, find the end of the object
							endOfObject := k + 1
							for l := k + 2; l < len(lines); l++ {
								line := strings.TrimSpace(lines[l])
								// Stop when we hit the next property at the same level as publishedListings
								if line != "" && !strings.HasPrefix(lines[l], "          ") && 
								   !strings.HasPrefix(lines[l], "            ") &&
								   strings.HasSuffix(line, ":") {
									endOfObject = l - 1
									break
								}
							}
							
							// Replace the entire items block with a reference
							indentation := strings.Repeat(" ", len(lines[k])-len(strings.TrimLeft(lines[k], " ")))
							newLines := make([]string, 0, len(lines))
							newLines = append(newLines, lines[:k+1]...)
							newLines = append(newLines, indentation+"  $ref: '#/components/schemas/PublishedListing'")
							newLines = append(newLines, lines[endOfObject+1:]...)
							lines = newLines
							changes++
							fmt.Printf("Fixed publishedListings inline schema to use PublishedListing schema reference\n")
							break
						}
					}
				}
			}
		}
	}

	// Add missing schema definitions to fix code-gen errors
	missingSchemas := []string{
		"PlatformHistory",
		"PublishedListing",
		"Otp",
	}
	
	// Find a good insertion point at the end of the schemas section
	schemasInsertIndex := -1
	inSchemasSection := false
	lastSchemaEndLine := -1
	
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		
		// Check if we're entering the schemas section
		if line == "schemas:" {
			inSchemasSection = true
			continue
		}
		
		// If we're in schemas section, track schema definitions
		if inSchemasSection {
			// Look for schema definitions (4 spaces + name + colon)
			if strings.HasPrefix(lines[i], "    ") && strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
				// This is a schema definition, find its end
				schemaIndentLevel := 4 // Schema level
				for j := i + 1; j < len(lines); j++ {
					nextLine := lines[j]
					if strings.TrimSpace(nextLine) == "" {
						continue // Skip empty lines
					}
					
					// Count leading spaces
					leadingSpaces := len(nextLine) - len(strings.TrimLeft(nextLine, " "))
					
					// If we hit another schema at the same level or a component section, we've found the end
					if leadingSpaces <= schemaIndentLevel && strings.TrimSpace(nextLine) != "" {
						if leadingSpaces < schemaIndentLevel || 
						   (leadingSpaces == schemaIndentLevel && strings.HasSuffix(strings.TrimSpace(nextLine), ":")) ||
						   strings.TrimSpace(nextLine) == "securitySchemes:" ||
						   strings.TrimSpace(nextLine) == "responses:" ||
						   strings.TrimSpace(nextLine) == "parameters:" {
							lastSchemaEndLine = j - 1
							if strings.TrimSpace(nextLine) == "securitySchemes:" ||
							   strings.TrimSpace(nextLine) == "responses:" ||
							   strings.TrimSpace(nextLine) == "parameters:" {
								schemasInsertIndex = j - 1
								goto found
							}
							break
						}
					}
				}
			}
		}
	}
	
	found:
	// If we didn't find a component section but found the end of the last schema, use that
	if schemasInsertIndex == -1 && lastSchemaEndLine > 0 {
		schemasInsertIndex = lastSchemaEndLine
	}
	
	// If we found a good insertion point, add the schemas
	if schemasInsertIndex > 0 {
		// Check if schemas already exist
		for _, schemaName := range missingSchemas {
			schemaExists := false
			for i := 0; i < len(lines); i++ {
				if strings.Contains(lines[i], "    "+schemaName+":") {
					schemaExists = true
					fmt.Printf("Schema %s already exists, skipping\n", schemaName)
					break
				}
			}
			
			if !schemaExists {
				var schemaLines []string
				switch schemaName {
				case "PlatformHistory":
					schemaLines = []string{
						"    PlatformHistory:",
						"      title: PlatformHistory",
						"      description: Platform History",
						"      type: object",
						"      properties:",
						"        isMobile:",
						"          type: boolean",
						"        platform:",
						"          $ref: '#/components/schemas/Platform'",
						"        recorded:",
						"          type: string",
						"          format: date-time",
					}
				case "PublishedListing":
					schemaLines = []string{
						"    PublishedListing:",
						"      title: PublishedListing", 
						"      description: Published Listing",
						"      type: object",
						"      properties:",
						"        description:",
						"          type: string",
						"        displayName:",
						"          type: string",
						"        imageId:",
						"          type: string",
						"        listingId:",
						"          type: string",
						"        listingType:",
						"          type: string",
						"        priceTokens:",
						"          type: integer",
					}
				case "Otp":
					schemaLines = []string{
						"    Otp:",
						"      title: Otp",
						"      description: Otp",
						"      type: object",
						"      required:",
						"        - code",
						"        - used",
						"      properties:",
						"        code:",
						"          type: string",
						"        used:",
						"          type: boolean",
					}
				}

				// Insert schema at the end of schemas section
				newLines := make([]string, len(lines)+len(schemaLines))
				copy(newLines[:schemasInsertIndex+1], lines[:schemasInsertIndex+1])
				copy(newLines[schemasInsertIndex+1:schemasInsertIndex+1+len(schemaLines)], schemaLines)
				copy(newLines[schemasInsertIndex+1+len(schemaLines):], lines[schemasInsertIndex+1:])
				lines = newLines
				schemasInsertIndex += len(schemaLines)
				changes++
				fmt.Printf("Added missing schema: %s\n", schemaName)
			}
		}
	}

	if changes == 0 {
		fmt.Println("No changes were made. The patterns might have already been fixed or not found.")
		return
	}

	// Join lines back together
	modifiedContent := strings.Join(lines, "\n")

	// Write the modified content back to the file
	err = os.WriteFile(filename, []byte(modifiedContent), 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Printf("Successfully applied %d schema fixes to %s\n", changes, filename)
}