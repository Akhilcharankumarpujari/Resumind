import { Link, useNavigate, useParams } from "react-router";
import { useEffect, useState } from "react";
import { useAuthStore } from "~/lib/auth";
import Summary from "~/components/Summary";
import ATS from "~/components/ATS";
import Details from "~/components/Details";

export const meta = () => ([
  { title: "Resumind | Resume Review" },
  { name: "description", content: "Detailed AI-powered overview of your resume" },
]);

interface ResumeData {
  id: string;
  companyName: string;
  jobTitle: string;
  jobDescription: string;
  feedback: any;
  createdAt: string;
}

const Resume = () => {
  const { isAuthenticated, isLoading } = useAuthStore();
  const { id } = useParams();
  const [pdfUrl, setPdfUrl] = useState("");
  const [feedback, setFeedback] = useState<any>(null);
  const [resumeData, setResumeData] = useState<ResumeData | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const navigate = useNavigate();


  useEffect(() => {
    if (!isLoading && !isAuthenticated) navigate(`/auth?next=/resume/${id}`);
  }, [isLoading, isAuthenticated, id, navigate]);


  useEffect(() => {
    if (!isAuthenticated || !id) return;

    const loadResume = async () => {
      try {

        const res = await fetch(`/api/resume/${id}`, {
          credentials: "include",
        });

        if (!res.ok) {
          setLoadError("Resume not found or you do not have access to it.");
          return;
        }

        const data: ResumeData = await res.json();
        setResumeData(data);
        setFeedback(data.feedback);



        setPdfUrl(`/api/resume/${id}/download`);
      } catch (err) {
        setLoadError("Failed to load resume data. Please try again.");
        console.error("Failed to load resume:", err);
      }
    };

    loadResume();
  }, [isAuthenticated, id]);

  return (
    <main className="!pt-0">
      <nav className="resume-nav">
        <Link to="/" className="back-button">
          <img src="/icons/back.svg" alt="logo" className="w-2.5 h-2.5" />
          <span className="text-gray-800 text-sm font-semibold">Back to Homepage</span>
        </Link>
      </nav>

      {loadError && (
        <div className="flex items-center justify-center h-64">
          <p className="text-red-500 text-center">{loadError}</p>
        </div>
      )}

      {!loadError && (
        <div className="flex flex-row w-full max-lg:flex-col-reverse">
          {}
          <section className="feedback-section bg-[url('/images/bg-small.svg')] bg-cover h-[100vh] sticky top-0 items-center justify-center">
            {pdfUrl ? (
              <div className="animate-in fade-in duration-1000 gradient-border max-sm:m-0 h-[90%] max-wxl:h-fit w-fit">
                <iframe
                  src={pdfUrl}
                  title="Resume PDF"
                  className="w-full h-full rounded-2xl"
                  style={{ minWidth: 320, minHeight: 400 }}
                />
              </div>
            ) : (
              <img
                src="/images/resume-scan-2.gif"
                className="w-[200px]"
                alt="Loading resume..."
              />
            )}
          </section>

          {}
          <section className="feedback-section">
            <h2 className="text-4xl !text-black font-bold">
              Resume Review
              {resumeData && (
                <span className="block text-lg font-normal text-gray-500 mt-1">
                  {resumeData.jobTitle} at {resumeData.companyName}
                </span>
              )}
            </h2>

            {feedback ? (
              <div className="flex flex-col gap-8 animate-in fade-in duration-1000">
                <Summary feedback={feedback} />
                <ATS score={feedback.ATS?.score || 0} suggestions={feedback.ATS?.tips || []} />
                <Details feedback={feedback} />
              </div>
            ) : (
              !loadError && (
                <img
                  src="/images/resume-scan-2.gif"
                  className="w-full"
                  alt="Loading feedback..."
                />
              )
            )}
          </section>
        </div>
      )}
    </main>
  );
};

export default Resume;
