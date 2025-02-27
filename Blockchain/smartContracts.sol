// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract PrescriptionNFT is ERC721, Ownable {
    uint256 private _tokenIdCounter;
    mapping(uint256 => Prescription) private _prescriptions;

    struct Prescription {
        string medication;
        string dosage;
        bool isActive;
    }

    event PrescriptionMinted(uint256 tokenId, address patient, string medication, string dosage);
    event PrescriptionDispensed(uint256 tokenId);

    constructor(address initialOwner) ERC721("PrescriptionNFT", "PRX") Ownable(initialOwner) {
        _tokenIdCounter = 1;
    }

    function mintPrescription(address patient, string memory medication, string memory dosage) external onlyOwner returns (uint256) {
        uint256 tokenId = _tokenIdCounter;
        _mint(patient, tokenId);

        _prescriptions[tokenId] = Prescription({
            medication: medication,
            dosage: dosage,
            isActive: true
        });

        _tokenIdCounter += 1;

        emit PrescriptionMinted(tokenId, patient, medication, dosage);

        return tokenId;
    }

    function dispensePrescription(uint256 tokenId) external {
        require(ownerOf(tokenId) == msg.sender, "Only the prescription holder can dispense it");
        require(_prescriptions[tokenId].isActive, "Prescription is not active");

        _prescriptions[tokenId].isActive = false;
        _burn(tokenId);

        emit PrescriptionDispensed(tokenId);
    }

    function getPrescriptionDetails(uint256 tokenId) external view returns (string memory medication, string memory dosage, bool isActive) {
        Prescription memory prescription = _prescriptions[tokenId];
        return (prescription.medication, prescription.dosage, prescription.isActive);
    }
}