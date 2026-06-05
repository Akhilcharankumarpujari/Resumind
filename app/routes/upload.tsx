import { type FormEvent, useState } from "react";
import Navbar from "~/components/Navbar";
import FileUploader from "~/components/FileUploader";
import { useNavigate } from "react-router";
import { useAuthStore } from "~/lib/auth";

export function meta() {
  return [
    { title: "Resumind | Upload Resume" },
    { name: "description", content: "Upload your resume for AI-powered feedback" },
  ];
}

const Upload = () => {
  const { isAuthenticated } = useAuthStore();
  const navigate = useNavigate();
  const [isProcessing, setIsProcessing] = useState(false);
  const [statusText, setStatusText] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleFileSelect = (selectedFile: File | null) => {
    setFile(selectedFile);
    setError(null);
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!file) {
      setError("Please select a PDF file to upload.");
      return;
    }

    if (!isAuthenticated) {
      navigate("/auth?next=/upload");
      return;
    }

    const form = e.currentTarget;
    const formData = new FormData(form);


    formData.set("resume", file);

    setIsProcessing(true);
    setError(null);
    setStatusText("Uploading your resume...");

    try {
      const res = await fetch("/api/analyze", {
        method: "POST",
        credentials: "include",
        body: formData,


      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.error || "Analysis failed");
      }

      setStatusText("Analyzing with AI...");
      const data = await res.json();

      setStatusText("Analysis complete! Redirecting...");
      navigate(`/resume/${data.id}`);
    } catch (err: any) {
      setError(err.message || "Something went wrong. Please try again.");
      setIsProcessing(false);
      setStatusText("");
    }
  };

  return (
    <main className="bg-[url('/images/bg-main.svg')] bg-cover">
      <Navbar />

      <section className="main-section">
        <div className="page-heading py-16">
          <h1>Smart feedback for your dream job</h1>

          {isProcessing ? (
            <>
              <h2>{statusText}</h2>
              <img src="/images/resume-scan.gif" className="w-full" alt="Analyzing..." />
            </>
          ) : (
            <h2>Drop your resume for an ATS score and improvement tips</h2>
          )}

          {}
          {error && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
              {error}
            </div>
          )}

          {!isProcessing && (
            <form
              id="upload-form"
              onSubmit={handleSubmit}
              className="flex flex-col gap-4 mt-8"
            >
              <div className="form-div">
                <label htmlFor="company-name">Company Name</label>
                <input
                  type="text"
                  name="companyName"
                  placeholder="e.g. Google"
                  id="company-name"
                  required
                />
              </div>

              <div className="form-div">
                <label htmlFor="job-title">Job Title</label>
                <input
                  type="text"
                  name="jobTitle"
                  placeholder="e.g. Software Engineer"
                  id="job-title"
                  required
                />
              </div>

              <div className="form-div">
                <label htmlFor="job-description">Job Description</label>
                <textarea
                  rows={5}
                  name="jobDescription"
                  placeholder="Paste the job description here for accurate AI analysis..."
                  id="job-description"
                />
              </div>

              <div className="form-div">
                <label htmlFor="uploader">Upload Resume (PDF only)</label>
                <FileUploader onFileSelect={handleFileSelect} />
                {file && (
                  <p className="text-sm text-green-600 mt-1">
                    ✅ Selected: {file.name}
                  </p>
                )}
              </div>

              <button className="primary-button" type="submit" disabled={!file}>
                Analyze Resume
              </button>
            </form>
          )}
        </div>
      </section>
    </main>
  );
};

export default Upload;
