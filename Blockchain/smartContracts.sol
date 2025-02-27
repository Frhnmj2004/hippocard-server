// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract PrescriptionNFT is ERC721, Ownable {
    uint256 private _tokenIdCounter;
    mapping(uint256 => Prescription) private _prescriptions;
    mapping(address => bool) private _doctors;
    mapping(address => bool) private _pharmacists;
    mapping(address => bool) private _hospitals;  // Emergency Access

    struct Prescription {
        string medication;
        string dosage;
        bool isActive;
    }

    event PrescriptionMinted(uint256 tokenId, address patient, string medication, string dosage);
    event PrescriptionDispensed(uint256 tokenId);

    modifier onlyDoctor() {
        require(_doctors[msg.sender], "Only doctors can perform this action");
        _;
    }

    modifier onlyPharmacist() {
        require(_pharmacists[msg.sender], "Only pharmacists can perform this action");
        _;
    }

    modifier onlyHospital() {
        require(_hospitals[msg.sender], "Only hospitals have emergency access");
        _;
    }

    constructor(address initialOwner) ERC721("PrescriptionNFT", "PRX") Ownable(initialOwner) {
        _tokenIdCounter = 1;
        _doctors[initialOwner] = true; // Contract owner is a doctor by default
    }

    // 🔹 **Doctor can mint a new prescription NFT**
    function mintPrescription(
        address patient,
        string memory medication,
        string memory dosage
    ) external onlyDoctor returns (uint256) {
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

    // 🔹 **View prescription details (based on role)**
    function getPrescriptionDetails(uint256 tokenId) 
        external 
        view 
        returns (string memory medication, string memory dosage, bool isActive) 
    {
        require(
            _doctors[msg.sender] || ownerOf(tokenId) == msg.sender || 
            (_pharmacists[msg.sender] && _prescriptions[tokenId].isActive) ||
            _hospitals[msg.sender], // Hospital emergency access
            "Access denied"
        );

        Prescription memory prescription = _prescriptions[tokenId];
        return (prescription.medication, prescription.dosage, prescription.isActive);
    }

    // 🔹 **View full medical history (Emergency Access for Hospitals)**
    function getAllPrescriptionsForPatient(address patient) 
        external 
        view 
        onlyHospital 
        returns (Prescription[] memory) 
    {
        uint256 totalCount = balanceOf(patient);
        Prescription[] memory history = new Prescription[](totalCount);

        uint256 counter = 0;
        for (uint256 tokenId = 1; tokenId < _tokenIdCounter; tokenId++) {
            if (_exists(tokenId) && ownerOf(tokenId) == patient) {
                history[counter] = _prescriptions[tokenId];
                counter++;
            }
        }
        return history;
    }

    // 🔹 **Pharmacist can mark a prescription as dispensed and burn NFT**
    function dispensePrescription(uint256 tokenId) external onlyPharmacist {
        require(_prescriptions[tokenId].isActive, "Prescription is already dispensed");

        _prescriptions[tokenId].isActive = false;
        _burn(tokenId);

        emit PrescriptionDispensed(tokenId);
    }

    // 🔹 **Admin functions to assign roles**
    function addDoctor(address doctor) external onlyOwner {
        _doctors[doctor] = true;
    }

    function removeDoctor(address doctor) external onlyOwner {
        _doctors[doctor] = false;
    }

    function addPharmacist(address pharmacist) external onlyOwner {
        _pharmacists[pharmacist] = true;
    }

    function removePharmacist(address pharmacist) external onlyOwner {
        _pharmacists[pharmacist] = false;
    }

    function addHospital(address hospital) external onlyOwner {
        _hospitals[hospital] = true;
    }

    function removeHospital(address hospital) external onlyOwner {
        _hospitals[hospital] = false;
    }
}
